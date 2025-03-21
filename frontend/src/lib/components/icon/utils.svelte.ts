import { EntityKind } from '$lib/api/canturin';
import BusIcon from './bus-icon.svelte';
import CanIdBuilderIcon from './can-id-builder-icon.svelte';
import MessageIcon from './message-icon.svelte';
import NetworkIcon from './network-icon.svelte';
import NodeIcon from './node-icon.svelte';
import SignalEnumIcon from './signal-enum-icon.svelte';
import SignalIcon from './signal-icon.svelte';
import SignalTypeIcon from './signal-type-icon.svelte';
import SignalUnitIcon from './signal-unit-icon.svelte';

export function getIconFromEntityKind(entKind: EntityKind) {
	switch (entKind) {
		case EntityKind.EntityKindNetwork:
			return NetworkIcon;
		case EntityKind.EntityKindBus:
			return BusIcon;
		case EntityKind.EntityKindNode:
			return NodeIcon;
		case EntityKind.EntityKindMessage:
			return MessageIcon;
		case EntityKind.EntityKindSignal:
			return SignalIcon;
		case EntityKind.EntityKindSignalType:
			return SignalTypeIcon;
		case EntityKind.EntityKindSignalUnit:
			return SignalUnitIcon;
		case EntityKind.EntityKindSignalEnum:
			return SignalEnumIcon;
		case EntityKind.EntityKindCANIDBuilder:
			return CanIdBuilderIcon;
		default:
			return NetworkIcon;
	}
}
