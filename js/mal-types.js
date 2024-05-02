export class TalList extends Array {}

export const isTalList = (a) => a instanceof Array;

export class TabSymbol extends String {}

export const isTabSymbol = (a) => a instanceof TabSymbol;

export class TalFunc extends Function {
	constructor(f, ast, params, env) {
		f.ast = ast;
		f.params = params;
		f.env = env;
		return Object.setPrototypeOf(f, new.target.prototype);
	}
}

export class TalMacro extends TalFunc {}

export class TabVar {
	constructor(value) {
		this.value = value;
	}
}

export class TabMap {
	constructor(f) {
		return Object.setPrototypeOf(f, new.target.prototype);
	}
}

export const isTabMap = (a) => a instanceof TabMap;

export class TabNil {}

export const isTabNil = (a) => a instanceof TabNil;

export const isTalFunc = (func) => TalFunc.prototype.isPrototypeOf(func);

export const isJsFunc = (func) => Function.prototype.isPrototypeOf(func);

export const isTalNumber = (a) => a.constructor === Number;