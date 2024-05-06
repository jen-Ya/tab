import { quasiquote } from './quasiquote.js';
import { Env } from './env.js';
import path from 'path';
import { TabNil, TabSymbol, TabMacro, TabMap, TabFunc, isJsFunc, isTabFunc } from './types.js';
import { Debugger } from './debugger.js';
import { inspect } from 'util';

const isMacroCall = (ast, env) => {
	return (
		ast instanceof Array &&
		ast.length > 0 &&
		ast[0] instanceof TabSymbol &&
		env.find(ast[0]) instanceof TabMacro
	);
};

const macroexpand = async(ast, env) => {
	while(isMacroCall(ast, env)) {
		const macro = env.get(ast[0]);
		// if first param is '.caller', provide caller function also
		const params = macro.params[0]?.valueOf() === '.caller' ?
			ast :
			ast.slice(1);
		ast = await macro(...params);
	}
	return ast;
};

const evalAst = async(ast, env, callstack) => {
	switch(ast.constructor) {
		case TabSymbol:
			switch(ast.valueOf()) {
				case '.filename':
					return ast.position?.filename;
				case '.dirname':
					return ast.position?.filename === 'repl' ?
						process.cwd() :
						path.dirname(ast.position?.filename);
				case '.cwd': {
					return process.cwd();
				}
				case '.env':
					return env;
				case '.tabhome':
					if(!process.env.TABHOME) {
						throw new Error('Environment variable TABHOME not set');
					}
					return process.env.TABHOME;
			}
			return env.get(ast);
		case Array: {
			const result = [];
			for(let i = 0; i < ast.length; i++) {
				result[i] = await EVAL(ast[i], env, callstack);
			}
			return result;
		}
		case TabMap: {
			const result = {};
			for(const key in ast) {
				result[key] = await EVAL(ast[key], env, callstack);
			}
			return new TabMap(result);
		}
		default:
			return ast;
	}
};

const letForm = async(ast, env, callstack) => {
	// let a "hallo"
	const v = await EVAL(ast[2], env, callstack);
	env.set(ast[1], v);
	ast = v;
	return { final: true, ast, env };
};

const evalForm = async(ast, env, callstack) => {
	const [_ast, _env = env] = await evalAst(ast.slice(1), env, callstack);
	ast = await EVAL(_ast, _env, [...callstack, ast[0]]);
	return { final: true, ast, env };
};

// Also maybe it would be enough to implement as an immediatly invoked anonymous function
const withForm = async(ast, env, callstack) => {
	env = Env(env);
	for(let i = 0; i < ast[1].length; i += 2) {
		env.set(
			ast[1][i],
			await EVAL(
				ast[1][i + 1],
				env,
				callstack,
			),
		);
	}
	ast = wrapDo(ast.slice(2));
	return { final: false, ast, env };
};

const doForm = async(ast, env, callstack) => {
	for(let i = 1; i < ast.length - 1; i++) {
		await EVAL(ast[i], env, callstack);
	}
	ast = ast[ast.length - 1];
	return { final: false, ast, env };
};

const ifForm = async(ast, env, callstack) => {
	// if (> x "1000") "gross" "klein"
	const evaledCond = await EVAL(ast[1], env, callstack);
	const cond = evaledCond !== false && !(evaledCond instanceof TabNil);
	if(cond) {
		ast = ast[2];
		return { final: false, ast, env };
	}
	else if(ast.length >= 4) {
		ast = ast[3];
		return { final: false, ast, env };
	} else {
		ast = new TabNil;
		return { final: true, ast, env };
	}
};

const quoteForm = async(ast, env) => {
	return { final: true, ast: ast[1], env };
};

const qqExpandForm = async(ast, env) => {
	ast = quasiquote(ast[1]);
	return { final: true, ast, env };
};

const qqForm = async(ast, env) => {
	ast = quasiquote(ast[1]);
	return { final: false, ast, env };
};

const wrapDo = (astList) => {
	if(astList.length === 0) {
		return new TabNil;
	}
	if(astList.length === 1) {
		return astList[0];
	}
	return [new TabSymbol('do'), ...astList];
};

const lambdaForm = async(ast, env) => {
	const [, params, ...fnAst] = ast;
	const body = wrapDo(fnAst);
	ast = new TabFunc(
		async(...args) => {
			const newEnv = Env(env, params, args);
			return await EVAL(body, newEnv);
		},
		body,
		params,
		env,
	);
	return { final: true, ast, env };
};

const macroForm = async(ast, env) => {
	const [, params, ...fnAst] = ast;
	const body = wrapDo(fnAst);
	ast = new TabMacro(
		async(...args) => {
			const newEnv = Env(env, params, args);
			return await EVAL(body, newEnv);
		},
		body,
		params,
		env,
	);
	return { final: true, ast, env };
};

const macroexpandForm = async(ast, env) => {
	ast = await macroexpand(ast[1], env);
	return { final: true, ast, env };
};

const tryForm = async(ast, env, callstack) => {
	try {
		return {
			final: true,
			ast: await EVAL(ast[1], env, callstack),
			env,
			callstack,
		};
	} catch (error) {
		const newEnv = Env(env);
		const malerror = error?.mal || error.message;
		newEnv.set(ast[2][1], malerror);
		return {
			final: false,
			ast: ast[2][2],
			env: newEnv,
		};
	}
};

const specialForms = {
	'let': letForm,
	'eval': evalForm,
	'with': withForm,
	'do': doForm,
	'if': ifForm,
	'q': quoteForm,
	'qqexpand': qqExpandForm,
	'qq': qqForm,
	'f': lambdaForm,
	'macrof': macroForm,
	'macroexpand': macroexpandForm,
	'try': tryForm,
};

const callFuncForm = async(ast, env, callstack) => {
	const [func, ...args] = await evalAst(ast, env, callstack);
	if(isTabFunc(func)) {
		callstack = [
			...callstack,
			ast[0],
		];
		ast = func.ast,
		env = Env(func.env, func.params, args);
		return { final: false, ast, env, callstack };
	}
	if(isJsFunc(func)) {
		ast = await func(...args);
		return { final: true, ast, env, callstack };
	}
	throw new Error(`Not a function: ${ func.constructor.name } (${ inspect(func) }) at ${ inspect(func.position) }`);
};

const evalError = (message, callstack) => {
	let prefix = '';
	for(const frame of callstack) {
		if(!frame.position) {
			prefix += `CallStack unknown position (${ frame })\n`;
			continue;
		}
		const {
			filename,
			start: [line, char],
		} = frame.position;
		prefix += `CallStack ${ filename }:${ line + 1 }:${ char + 1 }\n`;
	}
	const error = new Error(message);
	error.callstack = callstack;
	error.toString = () => prefix + message;
	return error;
};

export const EVAL = async(ast, env, callstack = []) => {
	// const { printString } = require('./printer');
	// console.log(printString(ast));
	try {
		for(;;) {
			await Debugger.evalStep(ast, env, callstack);
			if(!(ast instanceof Array)) {
				return await evalAst(ast, env, callstack);
			}
			if(ast.length === 0) {
				return ast;
			}
			ast = await macroexpand(ast, env);
			if(!(ast instanceof Array)) {
				return await evalAst(ast, env, callstack);
			}
			let final = false;
			const specialForm = specialForms[ast[0].valueOf()];
			if(specialForm) {
				({ final, ast, env } = await specialForm(ast, env, callstack));
			} else {
				({ final, ast, env, callstack } = await callFuncForm(ast, env, callstack));
			}
			if(final) {
				return ast;
			}
			continue;
		}
	} catch (error) {
		// TODO: check if caught
		if(Debugger.pauseOnEvalError) {
			await Debugger.didStop('error', ast, env, callstack);
		}
		// TODO: check if caught
		if(Debugger.dumpOnEvalError) {
			console.error('Error in evaluation');
			// console.error('ast', ast);
			if(ast.position) {
				console.error('ast position:', `${ ast.position.filename }:${ ast.position.start[0] + 1 }:${ ast.position.start[1] + 1 }`);
			}
			// console.error('env', env);
			// console.error('callstack', callstack);
		}
		throw error;
	}
};