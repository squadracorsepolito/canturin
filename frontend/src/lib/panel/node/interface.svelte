<script lang="ts">
	import { colorByName } from '$lib/actions/color-name.svelte';
	import type { BaseEntity, NodeInterface } from '$lib/api/canturin';
	import { Attribute } from '$lib/components/attribute';
	import { IconButton, LinkButton } from '$lib/components/button';
	import { Collapsible } from '$lib/components/collapsible';
	import Divider from '$lib/components/divider/divider.svelte';
	import { HoverPreview } from '$lib/components/hover-preview';
	import { AddIcon, DeleteIcon } from '$lib/components/icon';
	import { Readonly } from '$lib/components/readonly';
	import { Select } from '$lib/components/select';
	import { Table, TableField, TableTitle } from '$lib/components/table';
	import layout from '$lib/state/layout-state.svelte';
	import type { PanelSectionProps } from '../types';
	import { getNodeState } from './state.svelte';

	type Props = {
		buses: BaseEntity[];
		int: NodeInterface;
	} & PanelSectionProps;

	let { buses, int, entityId }: Props = $props();

	const ns = getNodeState(entityId);

	function handleAttachBus(busEntityID: string) {
		ns.updateAttachedBus(int.number, busEntityID);
	}

	let invalidBusIds = $derived.by(() => {
		if (!ns.entity.interfaces) return [];

		const ids = [];
		for (const int of ns.entity.interfaces) {
			if (int.number !== int.number) {
				ids.push(int.attachedBus.entityId);
			}
		}
		return ids;
	});

	function getCount(arr: any[] | null) {
		if (!arr) return 0;
		return arr.length;
	}

	function handleAddSentMessage() {
		ns.addSentMessage(int.number);
	}

	function handleBulkDeleteSentMessages(msgs: BaseEntity[]) {
		ns.removeSentMessages(
			int.number,
			msgs.map((m) => m.entityId)
		);
	}

	function handleDeleteSentMessage(msg: BaseEntity) {
		ns.removeSentMessage(int.number, msg.entityId);
	}

	function handleBulkDeleteReceivedMessages(msgs: BaseEntity[]) {
		ns.removeReceivedMessages(
			int.number,
			msgs.map((m) => m.entityId)
		);
	}

	function handleDeleteReceivedMessage(msg: BaseEntity) {
		ns.removeReceivedMessage(int.number, msg.entityId);
	}
</script>

{#snippet messagePreview(msg: BaseEntity)}
	<div class="font-medium text-sm pr-1">{msg.name}</div>

	{#if msg.desc}
		<div class="text-xs text-dimmed pt-1">{msg.desc}</div>
	{/if}
{/snippet}

<div class="flex gap-4">
	<div use:colorByName={{ name: int.attachedBus.name }} class="block w-2 rounded-btn"></div>

	<div class="flex-1">
		<div class="grid grid-cols-2 gap-5">
			<Attribute label="Number" desc="The number of the interface">
				<Readonly>
					<span class="font-medium">{int.number}</span>
				</Readonly>
			</Attribute>

			<Attribute label="Attached Bus" desc="The bus the interface is attached to">
				<Select
					items={buses}
					valueKey="entityId"
					name={`node-${int.number}-bus`}
					labelKey="name"
					bind:selected={int.attachedBus.entityId}
					filter={(item) => {
						return invalidBusIds.includes(item.entityId);
					}}
					onselect={handleAttachBus}
				/>
			</Attribute>
		</div>

		<Divider />

		<Collapsible initialCollapsed>
			{#snippet trigger()}
				<h4>Sent Messages ({getCount(int.sentMessages)})</h4>
			{/snippet}

			{#snippet content()}
				{#if int.sentMessages}
					<Table items={int.sentMessages} idKey="entityId">
						{#snippet bulkActions({ selectedCount, selectedItems, deselectAll })}
							<div class="flex gap-5">
								<IconButton
									onclick={() => {
										handleAddSentMessage();
										deselectAll();
									}}
									label="Add Sent Message"
									themeColor="primary"
								>
									<AddIcon />
								</IconButton>

								<IconButton
									onclick={() => handleBulkDeleteSentMessages(selectedItems)}
									themeColor="error"
									disabled={selectedCount === 0}
									label={`Delete Sent Messages ${selectedCount > 0 ? ` (${selectedCount})` : ''}`}
								>
									<DeleteIcon />
								</IconButton>
							</div>
						{/snippet}

						{#snippet header()}
							<TableTitle>Name</TableTitle>
						{/snippet}

						{#snippet row(msg)}
							<TableField>
								<HoverPreview placement="right">
									{#snippet trigger()}
										<LinkButton
											label={msg.name}
											onclick={() => layout.openMessagePanel(msg.entityId)}
										/>
									{/snippet}

									{#snippet content()}
										{@render messagePreview(msg)}
									{/snippet}
								</HoverPreview>
							</TableField>
						{/snippet}

						{#snippet rowActions(msg)}
							<IconButton onclick={() => handleDeleteSentMessage(msg)} themeColor="error">
								<DeleteIcon />
							</IconButton>
						{/snippet}
					</Table>
				{/if}
			{/snippet}
		</Collapsible>

		<Divider />

		<Collapsible initialCollapsed>
			{#snippet trigger()}
				<h4>Received Messages ({getCount(int.receivedMessages)})</h4>
			{/snippet}

			{#snippet content()}
				{#if int.receivedMessages}
					<Table items={int.receivedMessages} idKey="entityId">
						{#snippet bulkActions({ selectedCount, selectedItems })}
							<IconButton
								onclick={() => handleBulkDeleteReceivedMessages(selectedItems)}
								themeColor="error"
								disabled={selectedCount === 0}
								label={`Delete Received Messages ${selectedCount > 0 ? ` (${selectedCount})` : ''}`}
							>
								<DeleteIcon />
							</IconButton>
						{/snippet}

						{#snippet header()}
							<TableTitle>Name</TableTitle>
						{/snippet}

						{#snippet row(msg)}
							<TableField>
								<HoverPreview placement="right">
									{#snippet trigger()}
										<LinkButton
											label={msg.name}
											onclick={() => layout.openMessagePanel(msg.entityId)}
										/>
									{/snippet}

									{#snippet content()}
										{@render messagePreview(msg)}
									{/snippet}
								</HoverPreview>
							</TableField>
						{/snippet}

						{#snippet rowActions(msg)}
							<IconButton onclick={() => handleDeleteReceivedMessage(msg)} themeColor="error">
								<DeleteIcon />
							</IconButton>
						{/snippet}
					</Table>
				{/if}
			{/snippet}
		</Collapsible>
	</div>
</div>
