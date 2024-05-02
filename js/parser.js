export const parse = (tokens) => {
	let next = tokens[0];
	let rest = tokens;
	const parseError = (message) => {
		const position = next.position;
		if(!position) {
			return new Error(`ParseError at unkown location\n${ message }`);
		}
		return new Error(`ParserError at ${ position.filename }:${ position.start[0] + 1 }:${ position.start[1] + 1 }\n${ message }`);
	};
	const peek = (...expected) => {
		for(let i = 0; i < expected.length; i++) {
			if(rest[i].type !== expected[i]) {
				return false;
			}
		}
		return true;
	};
	const consume = (expected) => {
		if(!peek(expected)) {
			console.error('rest:', JSON.stringify(rest, null, '\t'));
			throw parseError(`unexpected token ${ next.type } ${ next.value }, expected ${ expected }`);
		}
		// console.log('consumed', expected, next.value)
		const consumed = next;
		rest = rest.slice(1);
		next = rest[0];
		return consumed;
	};
	const indent = () => consume('indent');
	const dedent = () => consume('dedent');
	const open = () => consume('(');
	const close = () => consume(')');
	const symbol = () => consume('symbol');
	const string = () => consume('string');
	const number = () => consume('number');
	const nil = () => consume('nil');
	const boolean = () => consume('boolean');

	const atom = () => {
		if(peek('symbol')) {
			return symbol();
		}
		if(peek('string')) {
			return string();
		}
		if(peek('number')) {
			return number();
		}
		if(peek('nil')) {
			return nil();
		}
		if(peek('boolean')) {
			return boolean();
		}
		throw parseError(`atom not implemented ${ next.type } / ${ next.value || '' }`);
	};
	const arg = () => {
		if(peek('(')) {
			const o = open();
			if(peek(')')) {
				const c = close();
				return {
					type: 'list',
					value: [],
					position: {
						filename: o.position.filename,
						start: o.position.start,
						end: c.position.end,
					},
				};
			}
			let exp;
			if(peek('indent')) {
				indent();
				exp = expression();
				dedent();
			} else {
				exp = expression();
			}
			const c = close();
			// TODO: Ugly noargs
			if(exp.type !== 'list' || exp.noargs) {
				return {
					type: 'list',
					value: [exp],
					position: {
						filename: o.position.filename,
						start: o.position.start,
						end: c.position.end,
					},
				};
			}
			return exp;
		}
		const _atom = atom();
		return _atom;
	};

	const inlineArgs = () => {
		const args = [];
		while(!peek('eol') && !peek(')') && !peek('eof') && !peek('dedent') && !peek('indent')) {
			const _arg = arg();
			args.push(_arg);
		}
		return args;
	};

	const indentArgs = () => {
		const args = [];
		if(peek('indent')) {
			consume('indent');
			for(;;) {
				if(peek('eol')) {
					consume('eol');
					continue;
				}
				if(peek('dedent')) {
					consume('dedent');
					break;
				}
				if(peek('eof')) {
					break;
				}
				const _arg = expression();
				args.push(_arg);
			}
		}
		return args;
	};

	const expression = () => {
		const first = arg();
		const args = [
			...inlineArgs(),
			...indentArgs(),
		];
		if(args.length === 0) {
			// TODO: Ugly noargs
			first.noargs = true;
			return first;
		}
		return {
			type: 'list',
			value: [
				first,
				...args,
			],
			position: {
				filename: first.position.filename,
				start: first.position.start,
				end: args[args.length - 1].position.end,
			},
		};
	};

	const lines = () => {
		const exps = [];
		while(!peek('eof')) {
			if(peek('eol')) {
				consume('eol');
				continue;
			}
			if(peek('comment')) {
				consume('comment');
				continue;
			}
			const exp = expression();
			exps.push(exp);
			while(peek('eol')) {
				consume('eol');
			}
		}
		if(exps.length === 1) {
			return exps[0];
		}
		// wrap multiple lines in do
		return {
			type: 'list',
			value: [
				{
					type: 'symbol',
					value: 'do',
				},
				...exps,
			],
			position: {
				filename: exps[0].position.filename,
				start: exps[0].position.start,
				end: exps[exps.length - 1].position.end,
			},
		};
	};
	return lines();
};