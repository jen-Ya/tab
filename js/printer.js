import { inspect } from 'util';
import { TabFunc, TabMap, TabVar, TabSymbol, TabNil } from './types.js';

const printList = (mal, start, stop, readable = true) => {
	return start + mal.map((val) => printString(val, readable)).join(' ') + stop;
};

const printHashmap = (mal, readable = true) => {
	const keys = Object.keys(mal);
	const printKey = key => '"' + key + '"' + ' ' + printString(mal[key], readable);
	return '(' + [
		'dict',
		...keys.map(printKey),
	].join(' ') + ')';
};

export const printString = (mal, readable = true) => {
	if(mal instanceof TabNil) {
		return 'nil';
	}
	if(mal === undefined) {
		return 'undefined';
	}
	if(
		[
			'boolean',
			'number',
		].indexOf(typeof mal) > -1
	) {
		return mal.toString();
	}
	if(typeof mal === 'string') {
		if(readable) {
			return '"' + mal
				.replace(/\\/g, '\\\\')
				.replace(/"/g, '\\"')
				.replace(/\n/g, '\\n')
				.replace(/\t/g, '\\t')
				+ '"';
		}
		return mal;
	}
	if(mal instanceof Array) {
		return printList(mal, '(', ')', readable);
	}
	if(mal instanceof TabMap) {
		return printHashmap(mal, readable);
	}
	if(mal instanceof TabFunc) {
		return '#<function>';
	}
	if(mal instanceof Function) {
		return '#<nativefunc>';
	}
	if(mal instanceof TabVar) {
		return `(var ${ printString(mal.value) })`;
	}
	if(mal instanceof TabSymbol) {
		return mal.valueOf();
	}
	return `#js<${ inspect(mal) }>`;
};