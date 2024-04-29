import { TalNil, TalMap } from './mal-types.js'

export const toJs = (mal) => {
	if(mal instanceof TalNil) {
		return null;
	}
	if(mal instanceof Array) {
		return mal.map(element => toJs(element));
	}
	if(mal instanceof TalMap) {
		let result = {};
		for(let key in mal) {
			result[key] = toJs(mal[key]);
		}
		return result;
	}
	return mal;
};

const simpleConstructors = new Set([
	Number,
	String,
	Boolean,
]);

export const fromJs = (js) => {
	if(js === null || js === undefined) {
		return new TalNil;
	}
	if(js instanceof Array) {
		return js.map(x => fromJs(x));
	}
	if(simpleConstructors.has(js.constructor)) {
		return js;
	}
	if(typeof js === 'object') {
		let result = {};
		for(let key in js) {
			result[key] = fromJs(js[key]);
		}
		return new TalMap(result);
	}
	return js;
};

export const callJs = async(func, args) => await func(...args);

const callMal = async(func, args) => {
	return await func(...args);
};