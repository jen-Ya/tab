import { Debugger } from './debugger.js';

export const output = (...args) => {
	/* eslint-disable no-console */
	console.log(...args);
	Debugger.output(...args);
};

export const outputError = (...args) => {
	console.error(...args);
	Debugger.output(...args);
};