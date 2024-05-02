import { readString as READ } from './reader.js';
import { printString as PRINT } from './printer.js';
import { EVAL } from './tab.js';
import { outputError } from './output.js';

export const rep = async(str, env) => {
	try {
		return PRINT(await EVAL(
			READ(
				str,
				{
					filename: 'repl',
				},
			),
			env,
		));
	} catch (error) {
		outputError(error);
		return error.message;
	}
};