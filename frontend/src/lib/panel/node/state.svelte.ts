import { BusService, NodeService, type Node } from '$lib/api/canturin';
import { HistoryNodeModify } from '$lib/constants/events';
import { pushToast } from '$lib/components/toast/toast-provider.svelte';
import { EntityState } from '$lib/state/entity-state.svelte';
import layout from '$lib/state/layout-state.svelte';
import { StateProvider } from '$lib/state/state-provider.svelte';

const provider = new StateProvider((node: Node) => new NodeState(node), HistoryNodeModify);

export function getNodeState(entityId: string) {
	return provider.get(entityId);
}

export async function loadNode(entityId: string) {
	const node = await NodeService.Get(entityId);
	provider.add(node);
}

export async function createNode(interfaceCount: number) {
	try {
		await NodeService.Create({ interfaceCount });
	} catch (err) {
		pushToast('error', 'Error', 'Operation failed');
		console.error(err);
	}
}

export async function deleteNode(entityId: string) {
	try {
		await NodeService.Delete(entityId);
		provider.remove(entityId);
		layout.closeIfOpen(entityId);
	} catch (error) {
		pushToast('error', 'Error', 'Operation failed');
		console.error(error);
	}
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
		const invalidIds = await NodeService.GetInvalidNodeIDs(this.entity.entityId);
		if (!invalidIds) return [];

		return invalidIds;
	}

	async getBuses() {
		const buses = await BusService.ListBase();
		if (!buses) return [];

		return buses;
	}

	updateName(name: string) {
		this.update(NodeService.UpdateName(this.entity.entityId, { name }));
	}

	updateDesc(desc: string) {
		this.update(NodeService.UpdateDesc(this.entity.entityId, { desc }));
	}

	updateID(id: number) {
		this.update(NodeService.UpdateNodeID(this.entity.entityId, { nodeId: id }));
	}

	updateAttachedBus(interfaceNumber: number, busEntityId: string) {
		this.update(
			NodeService.UpdateAttachedBus(this.entity.entityId, { interfaceNumber, busEntityId })
		);
	}

	addSentMessage(interfaceNumber: number) {
		this.update(NodeService.AddSentMessage(this.entity.entityId, { interfaceNumber }));
	}

	removeSentMessages(interfaceNumber: number, messageEntityIds: string[]) {
		this.update(
			NodeService.RemoveSentMessages(this.entity.entityId, {
				interfaceNumber,
				messageEntityIds
			})
		);
	}

	removeSentMessage(interfaceNumber: number, messageEntityId: string) {
		this.update(
			NodeService.RemoveSentMessages(this.entity.entityId, {
				interfaceNumber,
				messageEntityIds: [messageEntityId]
			})
		);
	}

	removeReceivedMessages(interfaceNumber: number, messageEntityIds: string[]) {
		this.update(
			NodeService.RemoveReceivedMessages(this.entity.entityId, {
				interfaceNumber,
				messageEntityIds
			})
		);
	}

	removeReceivedMessage(interfaceNumber: number, messageEntityId: string) {
		this.update(
			NodeService.RemoveReceivedMessages(this.entity.entityId, {
				interfaceNumber,
				messageEntityIds: [messageEntityId]
			})
		);
	}
}
