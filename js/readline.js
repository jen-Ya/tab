import { createInterface } from 'readline';
export const readline = createInterface({
	input: process.stdin,
	output: process.stdout,
});

export const ask = (question) => new Promise((resolve) => readline.question(question, resolve));