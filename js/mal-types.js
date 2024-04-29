export class TalList extends Array {}

export const isTalList = (a) => a instanceof Array;

export class TalSymbol extends String {}

export const isTalSymbol = (a) => a instanceof TalSymbol;

export class TalFunc extends Function {
	constructor(f, ast, params, env) {
		f.ast = ast;
		f.params = params;
		f.env = env;
		return Object.setPrototypeOf(f, new.target.prototype);
	}
}

export class TalMacro extends TalFunc {}

export class TalVar {
	constructor(value) {
		this.value = value;
	}
}

export class TalMap {
	constructor(f) {
		return Object.setPrototypeOf(f, new.target.prototype);
	}
}

export const isTalMap = (a) => a instanceof TalMap;

export class TalNil {}

export const isTalNil = (a) => a instanceof TalNil;

export const isTalFunc = (func) => TalFunc.prototype.isPrototypeOf(func);

export const isJsFunc = (func) => Function.prototype.isPrototypeOf(func);

export const isTalNumber = (a) => a.constructor === Number;