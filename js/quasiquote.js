import { TabSymbol, TabMap } from './mal-types.js';

const quasiquoteList = (ast) => {
	let result = [];
	for(let i = ast.length - 1; i >= 0; i--) {
		const elt = ast[i];
		if(
			elt instanceof Array &&
			elt.length > 0 &&
			elt[0] instanceof TabSymbol &&
			elt[0].valueOf() === '..unq'
		) {
			result = [
				new TabSymbol('concat'),
				elt[1],
				result,
			];
		} else {
			result = [
				new TabSymbol('cons'),
				quasiquote(elt),
				result,
			];
		}
	}
	return result;
};

export const quasiquote = (ast) => {
	if(
		ast instanceof Array &&
		ast.length > 0 &&
		ast[0] instanceof TabSymbol &&
		ast[0].valueOf() === 'unq'
	) {
		return ast[1];
	}
	if(ast instanceof Array) {
		return quasiquoteList(ast);
	}
	if(ast instanceof TabMap || ast instanceof TabSymbol) {
		return [
			new TabSymbol('quote'),
			ast,
		];
	}
	return ast;
};