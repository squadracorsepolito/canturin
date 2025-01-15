<script lang="ts">
	import { type Message, type Signal, SignalKind } from '$lib/api/canturin';
	import { IconButton, LinkButton } from '$lib/components/button';
	import { SignalGrid } from '$lib/components/grid';
	import { HoverPreview } from '$lib/components/hover-preview';
	import { CompactIcon, DeleteIcon } from '$lib/components/icon';
	import { Table, TableField, TableTitle } from '$lib/components/table';
	import layout from '$lib/state/layout-state.svelte';
	import type { PanelSectionProps } from '../types';
	import { getMessageState } from './state.svelte';
	import { getSignalKindString } from './utils';

	let { entityId }: PanelSectionProps = $props();

	const ms = getMessageState(entityId);

	function handleCompact() {
		ms.compactSignals();
	}

	function handleReorder(sigEntId: string, from: number, to: number) {
		ms.reorderSignal(sigEntId, from, to);
	}

	function handleBulkDelete(signals: Signal[]) {
		ms.deleteSignals(signals.map((sig) => sig.entityId));
	}

	function handleDelete(signal: Signal) {
		ms.deleteSignal(signal.entityId);
	}
</script>

{#snippet preview(sig: Signal)}
	<div>
		<span class="font-medium text-sm pr-1">{sig.name}</span>

		<span>
			{@render kindBadge(sig.kind)}
		</span>
	</div>

	{#if sig.desc}
		<div class="text-xs text-dimmed pt-1">{sig.desc}</div>
	{/if}
{/snippet}

{#snippet kindBadge(kind: SignalKind)}
	<span
		class={[
			'badge badge-sm',
			kind === SignalKind.SignalKindStandard && 'badge-primary',
			kind === SignalKind.SignalKindEnum && 'badge-secondary',
			kind === SignalKind.SignalKindMultiplexed && 'badge-accent'
		]}>{getSignalKindString(kind)}</span
	>
{/snippet}

{#snippet section(msg: Message)}
	{#if msg.signals}
		<div class="flex flex-col gap-5 @5xl:gap-2 @5xl:flex-row">
			<div class="flex-1">
				<Table items={msg.signals} idKey="entityId" reorder={handleReorder}>
					{#snippet bulkActions({ selectedCount, selectedItems, deselectAll })}
						<div class="flex justify-end gap-5">
							<IconButton onclick={handleCompact} color="secondary">
								<CompactIcon />
							</IconButton>

							<IconButton
								onclick={() => {
									handleBulkDelete(selectedItems);
									deselectAll();
								}}
								label={`Delete Signals ${selectedCount > 0 ? ` (${selectedCount})` : ''}`}
								disabled={selectedCount === 0}
								color="error"
							>
								<DeleteIcon />
							</IconButton>
						</div>
					{/snippet}

					{#snippet header()}
						<TableTitle>Name</TableTitle>

						<TableTitle>Kind</TableTitle>

						<TableTitle>Size</TableTitle>

						<TableTitle>Start Position</TableTitle>
					{/snippet}

					{#snippet row(sig)}
						<TableField>
							<HoverPreview placement="right">
								{#snippet trigger()}
									<LinkButton
										label={sig.name}
										onclick={() => layout.openPanel('signal', sig.entityId)}
									/>
								{/snippet}

								{#snippet content()}
									{@render preview(sig)}
								{/snippet}
							</HoverPreview>
						</TableField>

						<TableField>
							{@render kindBadge(sig.kind)}
						</TableField>

						<TableField>{sig.size}</TableField>

						<TableField>{sig.startPos}</TableField>
					{/snippet}

					{#snippet rowActions(signal)}
						<IconButton onclick={() => handleDelete(signal)} color="error">
							<DeleteIcon />
						</IconButton>
					{/snippet}
				</Table>
			</div>

			<div class="flex-1">
				<SignalGrid signals={msg.signals} height={msg.sizeByte} />
			</div>
		</div>
	{/if}
{/snippet}

<section>
	<h3 class="pb-5">Signals</h3>

	{@render section(ms.entity)}
</section>
