import { ns } from './core.js';
import { Env } from './env.js';
import { TabSymbol } from './mal-types.js';
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
	for(const key in ns) {
		env.set(new TabSymbol(key), ns[key]);
	}
	env.set(
		new TabSymbol('.env-root'),
		env,
	);
	await rep(initScript, env);
	const replEnv = Env(env);
	replEnv.set(
		new TabSymbol('.argv'),
		process.argv.slice(3),
	);
	return replEnv;
};