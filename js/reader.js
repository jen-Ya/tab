import { tokenize } from './tokenizer.js';
import { parse } from './parser.js';
import { TalSymbol, TalNil } from './mal-types.js';

const unpack = (parsed) => {
	let unpacked;
	if(parsed.type === 'list') {
		unpacked = parsed.value.map(element => unpack(element));
		unpacked.position = parsed.position;
	}
	else if(parsed.type === 'symbol') {
		unpacked = new TalSymbol(parsed.value);
		unpacked.position = parsed.position;
	}
	else if(parsed.type === 'string' || parsed.type === 'number' || parsed.type === 'boolean') {
		unpacked = parsed.value;
	}
	else if(parsed.type === 'nil') {
		unpacked = new TalNil;
		unpacked.position = parsed.position;
	}
	else {
		throw new Error(`Unexpected token type: ${ parsed.type }`);
	}
	return unpacked;
};

export const readString = (str, {
	filename = null,
} = {}) => {
	if(!str) {
		throw new Error('Unexpected EOF');
	}
	const tokens = tokenize(str, {
		filename,
	});
	return unpack(parse(tokens));
};