import { TalSymbol, TalMap } from './mal-types.js';

const quasiquoteList = (ast) => {
	let result = [];
	for(let i = ast.length - 1; i >= 0; i--) {
		const elt = ast[i];
		if(
			elt instanceof Array &&
			elt.length > 0 &&
			elt[0] instanceof TalSymbol &&
			elt[0].valueOf() === '..unq'
		) {
			result = [
				new TalSymbol('concat'),
				elt[1],
				result,
			];
		} else {
			result = [
				new TalSymbol('cons'),
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
		ast[0] instanceof TalSymbol &&
		ast[0].valueOf() === 'unq'
	) {
		return ast[1];
	}
	if(ast instanceof Array) {
		return quasiquoteList(ast);
	}
	if(ast instanceof TalMap || ast instanceof TalSymbol) {
		return [
			new TalSymbol('quote'),
			ast,
		];
	}
	return ast;
};