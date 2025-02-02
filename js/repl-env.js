import { ns } from './core.js';
import { Env } from './env.js';
import { TabSymbol } from './types.js';
import { rep } from './rep.js';

const initScript = `eval
	read-string
		file-read (str-join (li .tabhome '/init.tab') '')
		dict
			'filename'
			str-join (li .tabhome '/init.tab') ''
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