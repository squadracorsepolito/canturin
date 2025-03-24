import { CANIDBuilderOpKind } from '$lib/api/canturin';
import type { SegmentedControlOption } from '$lib/components/segmented-control/types';

export const canIdBuilderOpKindLabels = {
	[CANIDBuilderOpKind.$zero]: '',
	[CANIDBuilderOpKind.CANIDBuilderOpKindMessagePriority]: 'Message Priority',
	[CANIDBuilderOpKind.CANIDBuilderOpKindMessageID]: 'Message ID',
	[CANIDBuilderOpKind.CANIDBuilderOpKindNodeID]: 'Node ID',
	[CANIDBuilderOpKind.CANIDBuilderOpKindBitMask]: 'Bit Mask'
};

export const canIdBuilderOpKindDescs = {
	[CANIDBuilderOpKind.$zero]: '',
	[CANIDBuilderOpKind.CANIDBuilderOpKindMessagePriority]: 'Uses the message priority',
	[CANIDBuilderOpKind.CANIDBuilderOpKindMessageID]: 'Uses the message ID',
	[CANIDBuilderOpKind.CANIDBuilderOpKindNodeID]: 'Uses the node ID',
	[CANIDBuilderOpKind.CANIDBuilderOpKindBitMask]: 'Uses the bit mask'
};

export const canIdBuilderOpKindOptions: SegmentedControlOption[] = [
	{
		label: canIdBuilderOpKindLabels[CANIDBuilderOpKind.CANIDBuilderOpKindMessagePriority],
		value: CANIDBuilderOpKind.CANIDBuilderOpKindMessagePriority,
		desc: canIdBuilderOpKindDescs[CANIDBuilderOpKind.CANIDBuilderOpKindMessagePriority]
	},
	{
		label: canIdBuilderOpKindLabels[CANIDBuilderOpKind.CANIDBuilderOpKindMessageID],
		value: CANIDBuilderOpKind.CANIDBuilderOpKindMessageID,
		desc: canIdBuilderOpKindDescs[CANIDBuilderOpKind.CANIDBuilderOpKindMessageID]
	},
	{
		label: canIdBuilderOpKindLabels[CANIDBuilderOpKind.CANIDBuilderOpKindNodeID],
		value: CANIDBuilderOpKind.CANIDBuilderOpKindNodeID,
		desc: canIdBuilderOpKindDescs[CANIDBuilderOpKind.CANIDBuilderOpKindNodeID]
	},
	{
		label: canIdBuilderOpKindLabels[CANIDBuilderOpKind.CANIDBuilderOpKindBitMask],
		value: CANIDBuilderOpKind.CANIDBuilderOpKindBitMask,
		desc: canIdBuilderOpKindDescs[CANIDBuilderOpKind.CANIDBuilderOpKindBitMask]
	}
];
