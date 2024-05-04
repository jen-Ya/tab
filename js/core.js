import { printString } from './printer.js';
import { readString } from './reader.js';
import { readFile, writeFile, readdir } from 'fs/promises';
import {
	TabSymbol,
	TabVar,
	TabMap,
	TabNil,
	isTabNil,
	isTabList,
	isTabMap,
	isTabSymbol,
	isTabNumber,
} from './types.js';
import { output } from './output.js';
import { join as pathJoin, dirname, basename, resolve as pathResolve } from 'path';
import {
	fromJs,
	toJs,
	callJs,
} from './js-interop.js';
import { inspect } from 'util';
import { EVAL } from './tab.js';
import { Env } from './env.js';
import { tokenize } from './tokenizer.js';
import { parse } from './parser.js';
import { Debugger } from './debugger.js';

const eq = (a, b) => {
	if(a === undefined) {
		return b === undefined;
	}
	if(b === undefined) {
		return false;
	}
	if(
		Object.getPrototypeOf(a) !== Object.getPrototypeOf(b)
	) {
		return false;
	}
	if(isTabNil(a)) {
		return true;
	}
	if(isTabList(a)) {
		if(b.length !== a.length) {
			return false;
		}
		for(let i = 0; i < a.length; i++) {
			if(!eq(a[i], b[i])) {
				return false;
			}
		}
		return true;
	}
	if(isTabMap(a)) {
		const keys = Object.keys(a);
		if(keys.length !== Object.keys(b).length) {
			return false;
		}
		for(const key of keys) {
			if(!eq(a[key], b[key])) {
				return false;
			}
		}
		return true;
	}
	if(isTabSymbol(a)) {
		return a.valueOf() === b.valueOf();
	}
	return a === b;
};

const falsy = mal => isTabNil(mal) || mal === false;
const truthy = mal => !falsy(mal);

const _plus = (a, b) => a + b;
const _minus = (a, b) => a - b;
const _mul = (a, b) => a * b;
const _div = (a, b) => a / b;
const reducer = (func) => (start, ...args) => args.reduce(func, start);

const mathFuncs = {
	'+': reducer(_plus),
	'-': reducer(_minus),
	'*': reducer(_mul),
	'/': reducer(_div),
	'<': (a, b) => a < b,
	'<=': (a, b) => a <= b,
	'>': (a, b) => a > b,
	'>=': (a, b) => a >= b,
	'is-number': isTabNumber,
	'parse-number': (str) => {
		const num = parseFloat(str);
		if(!Number.isFinite(num)) {
			throw new Error(`parse-number: Invalid number [str=${ str }]`);
		}
		return num;
	},
};

const stringFuncs = {
	'char-at': (str, index) => str[index],
	'char-code': (char) => char.charCodeAt(0),
	'is-string': (arg) => arg.constructor === String,
	// todo: rename to str-slice, remove?
	'sub-str': (str, start, end) => str.slice(start, end),
	'str-len': (str) => str.length,
	// TODO: remove
	'str-ends-with': (str, suffix) => str.endsWith(suffix),
	// TODO: remove
	'str-starts-with': (str, prefix) => str.startsWith(prefix),
	// TODO: remove
	'str-join': (array, separator) => array.join(separator || ''),
	// TODO: remove
	'str-replace-all': (str, from, to) => str.replaceAll(from, to),
	'str-concat': (...args) => args.join(''),
};

const printFuncs = {
	'print-string': printString,
	output,
};

const listFuncs = {
	'List': Array,
	'li': (...args) => args,
	'count': (list) => isTabNil(list) ? 0 : list.length,
	'cons': (head, list) => [head, ...list],
	// TODO: is implemented in tab, but cannot be removed, since used in quasiquote
	'concat': (...lists) => lists.reduce((agg, arg) => [...agg, ...arg], []),
	'nth': (list, index) => {
		if(!isTabList(list)) {
			throw new Error(`nth argument list: Type error [expected=list|vector, got=${ list }]`);
		}
		if(index.constructor !== Number) {
			throw new Error(`nth argument index: Type error [expected=number, got=${ index.constructor.name }]`);
		}
		if(index >= list.length) {
			throw new Error(`nth argument index: Out of range [index=${ index }, length=${ list.length }] (${ inspect(list.position) })`);
		}
		return list[index];
	},
	'first': (list) => {
		if(list instanceof TabNil) {
			return new TabNil;
		}
		if(list.length === 0) {
			return new TabNil;
		}
		return list[0];
	},
	'slice': (array, start, stop) => array.slice(start, stop),
	'last': (array) => array[array.length - 1],
};

const typeFuncs = {
	'ist': (_type, arg) => Object.prototype.isPrototypeOf.call(_type.prototype, arg),
	'type': (arg) => arg.constructor.name,
};

const fileFuncs = {
	'file-read': (filename) => readFile(filename, 'utf-8'),
	// TODO: remove
	dirname,
	// TODO: remove
	basename,
	// TODO: remove
	'path-resolve': pathResolve,
	// TODO: remove
	'read-dir': readdir,
	// TODO: remove
	'path-join': pathJoin,
	// TODO: remove
	'file-write': async(filename, string) => {
		await writeFile(filename, string, 'utf-8');
		return new TabNil;
	},
};

const varFuncs = {
	'var': (mal) => new TabVar(mal),
	'Var': TabVar,
	'deref': (_var) => _var.value,
	'reset': (_var, mal) => {
		_var.value = mal;
		return _var;
	},
};

const funcFuncs = {
	'apply': async(...args) => {
		if(args.length < 2) {
			throw new Error(`Argument count error, [expected=2, got=${ args.length }`);
		}
		const func = args[0];
		const concats = args.slice(1, args.length - 1);
		let tail = args[args.length - 1];
		if(!(tail instanceof Array)) {
			tail = [tail];
		}
		return await func(
			...concats,
			...tail,
		);
	},
	Function,
};

const boolFuncs = {
	'not': falsy,
	'is-true': truthy,
	Boolean,
};

const dictFuncs = {
	'dict': (...args) => {
		const result = {};
		for(let i = 0; i < args.length; i += 2) {
			result[args[i]] = args[i + 1];
		}
		return new TabMap(result);
	},
	'Map': TabMap,
	'assoc': (hashmap, ...args) => {
		const result = {};
		for(const key in hashmap) {
			result[key] = hashmap[key];
		}
		for(let i = 0; i < args.length; i += 2) {
			result[args[i]] = args[i + 1];
		}
		return new TabMap(result);
	},
	'dissoc': (hashmap, ...keys) => {
		const result = {};
		const delset = new Set(keys);
		for(const key in hashmap) {
			if(!delset.has(key)) {
				result[key] = hashmap[key];
			}
		}
		return new TabMap(result);
	},
	'has': (hashmap, key) => {
		return key in hashmap;
	},
	'get': (hashmap, key) => {
		if(hashmap instanceof TabNil) {
			return new TabNil;
		}
		return hashmap[key] ?? new TabNil;
	},
	// TODO: Should dicts be mutable?
	'set': (hashmap, key, value) => {
		return new TabMap({
			...hashmap,
			[key]: value,
		});
	},
	'keys': Object.keys,
	'vals': Object.values,
	'entries': Object.entries,
};

const symbolFuncs = {
	'Symbol': TabSymbol,
	'symbol': (arg) => new TabSymbol(arg),
	'symbol-value': (arg) => arg.valueOf(),
};

const jsFuncs = {
	'call-js': async(func, ...args) => await callJs(func, args),
	'to-js': (mal) => toJs(mal),
	'from-js': (js) => fromJs(js),
	'js-eval': (string) => fromJs(eval(string)),
	'js-raw': (string) => eval(string),
	/* TODO: remove */
	'is-undefined': (arg) => arg === undefined,
};

const envFuncs = {
	'env-new': Env,
	'env-set': (env, key, value) => {
		env.set(key, value);
		return value;
	},
	'env-find': (env, key) => env.find(key),
	'env-get': (env, key) => env.get(key),
	'env-outer': (env) => env.outer || new TabNil,
};

const metaFuncs = {
	'read-string': (malstr, options) => readString(malstr, {
		filename: options?.['filename'] || null,
	}),
	tokenize,
	parse,
	EVAL,
	Debugger,
};

const eqFuncs = {
	'=': eq,
	'is': (a, b) => a === b,
};

const errorFuncs = {
	'throw': (mal) => {
		const error = new Error(mal.constructor === String ? mal : 'Error');
		error.mal = mal;
		throw error;
	},
};

const nilFuncs = {
	'Nil': TabNil,
};

const processFuncs = {
	'exit': (code) => process.exit(code),
};

export const ns = {
	...mathFuncs,
	...stringFuncs,
	...printFuncs,
	...listFuncs,
	...typeFuncs,
	...fileFuncs,
	...varFuncs,
	...funcFuncs,
	...boolFuncs,
	...dictFuncs,
	...symbolFuncs,
	...jsFuncs,
	...envFuncs,
	...metaFuncs,
	...eqFuncs,
	...errorFuncs,
	...nilFuncs,
	...processFuncs,
};