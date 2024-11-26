import { BusService, NodeService, type Node } from '$lib/api/canturin';
import { HistoryNodeModify } from '$lib/api/events';
import { EntityState } from '$lib/state/entity-state.svelte';
import { StateProvider } from '$lib/state/state-provider.svelte';

const provider = new StateProvider((node: Node) => new NodeState(node), HistoryNodeModify);

export function getNodeState(entityId: string) {
	return provider.get(entityId);
}

export async function loadNode(entityId: string) {
	const node = await NodeService.Get(entityId);
	provider.add(node);
}

class NodeState extends EntityState<Node> {
	constructor(node: Node) {
		super(node);
	}

	async getInvalidNames() {
		const invalidNames = await NodeService.GetInvalidNames(this.entity.entityId);

		if (invalidNames) {
			return invalidNames;
		}

		return [];
	}

	async getInvalidIds() {
		const invalidIds = await NodeService.GetInvalidIDs(this.entity.entityId);
		if (!invalidIds) return [];

		return invalidIds;
	}

	async getBuses() {
		const buses = await BusService.ListBase();
		if (!buses) return [];

		return buses;
	}

	updateName(name: string) {
		this.update(NodeService.UpdateName(this.entity.entityId, name));
	}

	updateDesc(desc: string) {
		this.update(NodeService.UpdateDesc(this.entity.entityId, desc));
	}

	updateID(id: number) {
		this.update(NodeService.UpdateID(this.entity.entityId, id));
	}

	attachBus(interfaceNumber: number, busEntityID: string) {
		this.update(NodeService.AttachBus(this.entity.entityId, interfaceNumber, busEntityID));
	}
}
