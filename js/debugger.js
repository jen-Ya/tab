import { once, EventEmitter } from 'events';

const compareCallStack = (callstackA, callstackB) => {
	if(callstackA === callstackB) {
		return true;
	}
	if(callstackA.length !== callstackB.length) {
		return false;
	}
	for(let i = 0; i < callstackA.length; i++) {
		if(callstackA[i] !== callstackB[i]) {
			return false;
		}
	}
	return true;
};

export const Debugger = {
	events: new EventEmitter(),
	shouldPause: null,
	pauseCallstack: null,
	ast: null,
	env: null,
	callstack: null,
	dumpOnEvalError: false,
	pauseOnEvalError: false,
	breakpoints: [],
	terminate: () => {
		const reason = 'stop';
		Debugger.shouldPause = reason;
		Debugger.events.emit('do-continue', reason);
	},
	setBreakpoints: (path, lines) => {
		Debugger.breakpoints = Debugger.breakpoints.filter(bp => bp.path !== path).concat(lines.map(line => ({
			path,
			line,
		})));
	},
	hasBreakpoint: (path, line) => {
		return Debugger.breakpoints.findIndex(bp => bp.path === path && bp.line === line) > -1;
	},
	output: (...args) => {
		Debugger.events.emit('did-output', ...args);
	},
	outputError: (...args) => {
		Debugger.events.emit('did-output', ...args);
	},
	shouldExaminePause: (ast) => {
		return Debugger.shouldPause !== null &&
			ast.position?.filename &&
			ast.position.filename !== 'repl' &&
			!ast.position.filename.endsWith('/init.tab');
	},
	shouldStepIn: (callstack) => {
		return Debugger.shouldPause === 'step-in' &&
			callstack.length >= 1 &&
			compareCallStack(
				Debugger.pauseCallstack,
				callstack.slice(0, Debugger.pauseCallstack.length),
			);
	},
	shouldStepOver: (callstack) => {
		return Debugger.shouldPause === 'step-over' &&
			compareCallStack(Debugger.pauseCallstack, callstack);
	},
	shouldStepOut: (callstack) => {
		return Debugger.shouldPause === 'step-out' &&
			compareCallStack(
				Debugger.pauseCallstack.slice(0, callstack.length),
				callstack,
			);
	},
	shouldPauseOnBreakpoint: (ast) => {
		return ast?.position?.filename &&
			Debugger.hasBreakpoint(ast.position.filename, ast.position.start[0]);
	},
	shouldInterrupt: (ast) => {
		return ast?.position?.filename && Debugger.shouldPause === 'interrupt';
	},
	didStop: async(reason, ast, env, callstack) => {
		Debugger.ast = ast;
		Debugger.env = env;
		Debugger.callstack = callstack;
		Debugger.events.emit('did-stop', ast, reason);
		const [pauseType] = await once(Debugger.events, 'do-continue');
		Debugger.shouldPause = pauseType;
		if(pauseType === 'step-in') {
			Debugger.pauseCallstack = callstack;
		}
		else if(pauseType === 'step-out') {
			if(callstack.length === 0) {
				Debugger.pauseCallstack = null;
			} else {
				Debugger.pauseCallstack = callstack.slice(0, callstack.length - 1);
			}
		}
		else if(pauseType === 'step-over') {
			Debugger.pauseCallstack = callstack;
		}
		else if(pauseType === 'stop') {
			throw new Error('stopped from debugger');
		}
	},
	shouldPauseOnStep: (ast, callstack) => {
		return Debugger.shouldExaminePause(ast) &&
			(
				Debugger.shouldStepIn(callstack) ||
				Debugger.shouldStepOver(callstack) ||
				Debugger.shouldStepOut(callstack)
			);
	},
	exitIfStopped: () => {
		if(Debugger.shouldPause === 'stop') {
			throw new Error('stopped from debugger');
		}
	},
	checkPause: (ast, callstack) => {
		if(Debugger.shouldPauseOnStep(ast, callstack)) {
			return 'step';
		}
		else if(Debugger.shouldPauseOnBreakpoint(ast)) {
			return 'breakpoint';
		}
		else if(Debugger.shouldInterrupt(ast)) {
			return 'pause';
		}
		return null;
	},
	pauseIfNeeded: async(ast, env, callstack) => {
		const pauseReason = Debugger.checkPause(ast, callstack);
		if(pauseReason) {
			await Debugger.didStop(pauseReason, ast, env, callstack);
		}
	},
	evalStep: async(ast, env, callstack) => {
		Debugger.exitIfStopped();
		await Debugger.pauseIfNeeded(ast, env, callstack);
	},
};