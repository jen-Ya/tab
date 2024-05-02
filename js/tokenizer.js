const digits = '1234567890';
const noid = '\n#\'"() \t';

const stringTransform = str => str
	.replace(/\\n/g, '\n')
	.replace(/\\t/g, '\t');

export const tokenize = (text, {
	keepComments = false,
	filename = null,
} = {}) => {
	text = text += '\n';
	let cursor = 0;
	let linecount = 0;
	let charcount = 0;
	let indentation = 0;
	const tokens = [];
	let tokenstart = [0, 0];

	const tokenizeError = (message) => {
		return new Error(`TokenizerError at ${ filename }:${ linecount + 1 }:${ charcount + 1 }\n` +
			`Token Started at ${ filename }:${ tokenstart[0] + 1 }:${ tokenstart[1] + 1 }\n` +
			message);
	};

	const inccursor = () => {
		if(cursor > text.length) {
			throw tokenizeError('TokenizerError: cursor out of bounds');
		}
		if(text[cursor] === '\n') {
			linecount++;
			charcount = 0;
		} else {
			charcount++;
		}
		cursor++;
	};

	const addtoken = (type, value) => {
		tokens.push({
			type,
			value,
			position: {
				filename,
				start: tokenstart,
				end: [linecount, charcount],
			},
		});
		tokenstart = [linecount, charcount];
	};

	const consume = (chars) => {
		for(const char of chars) {
			if(text[cursor] !== char) {
				throw tokenizeError(`TokenizeError: unexpected ${ text[cursor] }, expected ${ char }`);
			}
			inccursor();
		}
	};

	const isPeek = (chars) => {
		return text.slice(cursor, cursor + chars.length) === chars;
	};

	const isNoId = (char) => {
		return noid.indexOf(char) > -1;
	};

	const isPeekWord = (chars) => {
		return isPeek(chars) && isNoId(text[cursor + chars.length]);
	};

	const isDigit = (char) => {
		return digits.indexOf(char) > -1;
	};

	const consumeIndent = (indent) => {
		for(let i = 0; i < indent; i++) {
			consume('\t');
		}
	};

	const getIndent = () => {
		let indent = 0;
		while(indent < text.length - cursor) {
			if(text[cursor + indent] !== '\t') {
				return indent;
			}
			indent++;
		}
		return indent;
	};

	const eol = () => {
		while(text[cursor] === '\n') {
			consume('\n');
		}
		const indent = getIndent();
		if(indent === indentation) {
			consumeIndent(indent);
			addtoken('eol');
		}
		else if(indent < indentation) {
			consumeIndent(indent);
			for(let i = indent; i < indentation; i++) {
				addtoken('dedent');
			}
			indentation = indent;
		}
		else if(indent === indentation + 1) {
			consumeIndent(indent);
			addtoken('indent');
			indentation = indent;
		}
		else {
			throw tokenizeError(`IndentError: unexpected indentation of ${ indent } from ${ indentation }`);
		}
	};

	const space = () => consume(' ');

	const consumeMultiline = (separator) => {
		consume(separator);
		consume('\n');
		let value = '';
		consumeIndent(indentation + 1);
		for(;;) {
			while(text[cursor] !== '\n') {
				value += text[cursor];
				inccursor();
			}
			consume('\n');
			while(text[cursor] === '\n') {
				consume('\n');
				value += '\n';
			}
			const indent = getIndent();
			if(indent === indentation) {
				consumeIndent(indentation);
				consume(separator);
				return stringTransform(value);
			}
			value += '\n';
			consumeIndent(indentation + 1);
		}
	};

	const consumeToEol = () => {
		let value = '';
		while(text[cursor] !== '\n') {
			value += text[cursor];
			inccursor();
		}
		return value;
	};

	const comment = () => {
		if(text[cursor + 1] === '\n') {
			const value = consumeMultiline('#');
			if(keepComments) {
				addtoken('comment', value);
			}
		} else {
			consume('#');
			if(text[cursor] === ' ') {
				consume(' ');
			}
			const value = consumeToEol();
			if(keepComments) {
				addtoken('comment', value);
			}
		}
	};

	const string = (quote) => {
		if(text[cursor + 1] === '\n') {
			addtoken('string', consumeMultiline(quote));
		} else {
			consume(quote);
			let value = '';
			while(!isPeek(quote)) {
				if(isPeek('\n')) {
					throw tokenizeError(`Unterminated single line string ${ quote }`);
				}
				{ value += text[cursor]; }
				inccursor();
			}
			consume(quote);
			// Should we unescape escaped quotes in single line strings?
			// e.g. "foo\"bar" -> foo"bar?
			addtoken('string', stringTransform(value));
		}
	};

	const number = () => {
		let value = '';
		if(isPeek('-')) {
			value += text[cursor];
			consume('-');
		}
		while(isDigit(text[cursor])) {
			value += text[cursor];
			inccursor();
		}
		if(text[cursor] === '.') {
			consume('.');
			value += '.';
		}
		while(isDigit(text[cursor])) {
			value += text[cursor];
			inccursor();
		}
		addtoken('number', +value);
	};

	const eof = () => {
		addtoken('eof');
	};

	const symbol = () => {
		let value = '';
		while(!isNoId(text[cursor])) {
			value += text[cursor];
			inccursor();
		}
		if(!value) {
			throw tokenizeError('symbol has no value');
		}
		addtoken('symbol', value);
	};

	const open = () => {
		consume('(');
		addtoken('(');
	};

	const close = () => {
		consume(')');
		addtoken(')');
	};

	for(;;) {
		if(isPeek(' ')) {
			space();
		}
		else if(isPeek('#')) {
			comment();
		}
		else if(isPeek("'") || isPeek('"')) {
			string(text[cursor]);
		}
		else if(
			isDigit(text[cursor]) ||
			isPeek('-') && isDigit(text[cursor + 1])
		) {
			number();
		}
		else if(isPeek('(')) {
			open();
		}
		else if(isPeek(')')) {
			close();
		}
		else if(isPeek('\n')) {
			eol();
		}
		else if(isPeekWord('nil')) {
			consume('nil');
			addtoken('nil', null);
		}
		else if(isPeekWord('_')) {
			consume('_');
			addtoken('nil', null);
		}
		else if(isPeekWord('true')) {
			consume('true');
			addtoken('boolean', true);
		}
		else if(isPeekWord('false')) {
			consume('false');
			addtoken('boolean', false);
		}
		else if(cursor >= text.length) {
			eof();
			return tokens;
		}
		else {
			symbol();
		}
	}
};