import randomColor from 'randomcolor';

import { SignalKind } from './api/github.com/squadracorsepolito/acmelib';

export function genSignalColor(signalName: string) {
	return randomColor({
		seed: signalName
	});
}

export function getSignalKindName(sigKind: SignalKind): string {
	console.log('OK');

	switch (sigKind) {
		case SignalKind.SignalKindStandard:
			return 'Standard';
		case SignalKind.SignalKindEnum:
			return 'Enum';
		case SignalKind.SignalKindMultiplexer:
			return 'Multiplexer';
	}
}
