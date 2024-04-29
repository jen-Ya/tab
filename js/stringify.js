const toString = (ast) => {
	if(Array.isArray(ast)) {
		return ast.map((node) => toString(node)).join('\n');
	}
	switch(ast.type) {
		case 'list':
			return listToString(ast);
		case 'string':
			return '"' + ast.value + '"';
	}
	return ast.value;
};

const listToString = (list) => {
	return '(' + list.value.map(ast => toString(ast)).join(' ') + ')';
};

// const listToString = ({ value: [{ value: first}, ...args]}, indent = '') => {
// 	return `${ indent }${ first }\n${ args.map((arg) => toString(arg, indent + '  ')).join('\n') }`
// }

const tokensToString = tokens => tokens.map(({ type, value }) => `${ value || '' }:${ type }`).join(' ');

const tokensToStringFormatted = tokens => {
	let indent = 0;
	let output = '';
	for(let token of tokens) {
		switch(token.type) {
			case 'eol':
				output += '\n';
				break;
			case 'indent':
				indent++;
				output += '\n';
				break;
			case 'dedent':
				indent--;
				output += '\n';
				break;
			default:
				for(let i = 0; i < indent; i++) {
					output += '  ';
				}
				output += `${ token.value || '' }:${ token.type } `;
		}
	}
	return output;
};

module.exports = {
	toString,
	listToString,
	tokensToString,
	tokensToStringFormatted,
};