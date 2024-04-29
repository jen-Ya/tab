import { ns } from './core.js';
import { Env } from './env.js';
import { TalSymbol } from './mal-types.js';
import { rep } from './rep.js';
import path from 'path';

const initScript = `eval
	read-string
		file-read (str-join (list .tabhome '/init.tab') '')
		dict
			'filename'
			str-join (list .tabhome '/init.tab') ''
`;

export const makeReplEnv = async() => {
	const env = Env(null);
	for(let key in ns) {
		env.set(new TalSymbol(key), ns[key]);
	}
	env.set(
		new TalSymbol('.env-root'),
		env,
	);
	await rep(initScript, env);
	const replEnv = Env(env);
	replEnv.set(
		new TalSymbol('.argv'),
		process.argv.slice(3),
	);
	return replEnv;
};