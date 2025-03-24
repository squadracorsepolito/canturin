import { CANIDBuilderService, MessagePriority, type CANIDBuilder } from '$lib/api/canturin';
import { HistoryCANIDBuilderModify } from '$lib/constants/events';
import { EntityState } from '$lib/state/entity-state.svelte';
import { StateProvider } from '$lib/state/state-provider.svelte';

const provider = new StateProvider(
	(canIdBuilder: CANIDBuilder) => new CanIdBuilderState(canIdBuilder),
	HistoryCANIDBuilderModify
);

export function getCanIdBuilderState(entityId: string) {
	return provider.get(entityId);
}

export async function loadCanIdBuilder(entityId: string) {
	const canIdBuilder = await CANIDBuilderService.Get(entityId);
	return provider.add(canIdBuilder);
}

export class CanIdBuilderState extends EntityState<CANIDBuilder> {
	constructor(canIdBuilder: CANIDBuilder) {
		super(canIdBuilder);
	}

	async getInvalidNames() {
		const invalidNames = await CANIDBuilderService.GetInvalidNames(this.entity.entityId);

		if (invalidNames) {
			return invalidNames;
		}

		return [];
	}

	updateName(name: string) {
		this.update(CANIDBuilderService.UpdateName(this.entity.entityId, { name }));
	}

	updateDesc(desc: string) {
		this.update(CANIDBuilderService.UpdateDesc(this.entity.entityId, { desc }));
	}

	async calculateCANIds(messagePriority: MessagePriority, messageId: number, nodeId: number) {
		const canIds = await CANIDBuilderService.CalculateCANIDs(this.entity.entityId, {
			messagePriority,
			messageId,
			nodeId
		});

		if (canIds) return canIds;

		return [];
	}
}
