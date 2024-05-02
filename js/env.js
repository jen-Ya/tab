import { TabNil } from './types.js';

export const Env = (outer = null, binds = [], exprs = []) => {
	const data = {};
	const set = (symbol, value) => {
		data[symbol.valueOf()] = value;
	};
	if(!(binds instanceof TabNil)) {
		if(!(binds instanceof Array)) {
			binds = [binds];
		}
		for(let i = 0; i < binds.length; i++) {
			if(binds[i].valueOf() === '..') {
				set(binds[i + 1], exprs.slice(i));
				break;
			}
			set(binds[i], exprs[i]);
		}
	}
	const find = (symbol) => {
		let _data = data;
		let _outer = outer;
		const key = symbol.valueOf();
		for(;;) {
			if(key in _data) {
				return _data[key];
			}
			if(_outer === null) {
				return null;
			}
			_data = _outer.data;
			_outer = _outer.outer;
		}
	};
	const get = (symbol) => {
		const v = find(symbol);
		if(v === null) {
			throw new Error(`Cannot find symbol: ${ symbol }`);
		}
		if(v === undefined) {
			return new TabNil;
		}
		return v;
	};
	return {
		outer,
		data,
		binds,
		exprs,
		set,
		find,
		get,
	};
};