import { darken, getLuminance, lighten } from 'color2k';
import randomColor from 'randomcolor';

import type { ReferenceTreeNode } from './components/tree/types';
import type { SignalReference } from './api/canturin';

export function getHexNumber(num: number) {
	return `0x${num.toString(16).padStart(2, '0')}`;
}

export function getColorByName(name: string) {
	const bgColor = randomColor({
		seed: name
	});

	let textColor = '';
	if (getLuminance(bgColor) > 0.4) {
		textColor = darken(bgColor, 0.7);
	} else {
		textColor = lighten(bgColor, 0.7);
	}

	return {
		bgColor: bgColor,
		textColor: textColor
	};
}

export function getSignalReferenceTree(sigRefs: SignalReference[]): ReferenceTreeNode[] {
	const res: ReferenceTreeNode[] = [];

	for (const ref of sigRefs) {
		const signalNode: ReferenceTreeNode = {
			name: ref.signal.name,
			childNodes: []
		};

		const msgNode: ReferenceTreeNode = {
			name: ref.message.name,
			childNodes: [signalNode]
		};

		const nodeNode: ReferenceTreeNode = {
			name: ref.node.name,
			childNodes: [msgNode]
		};

		const busNode: ReferenceTreeNode = {
			name: ref.bus.name,
			childNodes: [nodeNode]
		};

		const bus = res.find((b) => b.name === ref.bus.name);
		if (bus === undefined) {
			res.push(busNode);
			continue;
		}

		const node = bus.childNodes.find((n) => n.name === ref.node.name);
		if (node === undefined) {
			bus.childNodes.push(nodeNode);
			continue;
		}

		const msg = node.childNodes.find((m) => m.name === ref.message.name);
		if (msg === undefined) {
			node.childNodes.push(msgNode);
			continue;
		}

		msg.childNodes.push(signalNode);
	}

	return res;
}

export function uniqueId() {
	return Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15);
}

export type KeyOfString<T> = {
	[K in keyof T]: T[K] extends string ? K : never;
}[keyof T];
