import { TabNil, TabMap } from './mal-types.js';

export const toJs = (mal) => {
	if(mal instanceof TabNil) {
		return null;
	}
	if(mal instanceof Array) {
		return mal.map(element => toJs(element));
	}
	if(mal instanceof TabMap) {
		const result = {};
		for(const key in mal) {
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
		return new TabNil;
	}
	if(js instanceof Array) {
		return js.map(x => fromJs(x));
	}
	if(simpleConstructors.has(js.constructor)) {
		return js;
	}
	if(typeof js === 'object') {
		const result = {};
		for(const key in js) {
			result[key] = fromJs(js[key]);
		}
		return new TabMap(result);
	}
	return js;
};

export const callJs = async(func, args) => await func(...args);

const callMal = async(func, args) => {
	return await func(...args);
};