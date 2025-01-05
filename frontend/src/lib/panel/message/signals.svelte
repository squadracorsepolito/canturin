<script lang="ts">
	import type { Message, Signal } from '$lib/api/canturin';
	import { IconButton } from '$lib/components/button';
	import { SignalGrid } from '$lib/components/grid';
	import { CompactIcon, DeleteIcon } from '$lib/components/icon';
	import { Table, TableField, TableTitle } from '$lib/components/table';
	import type { PanelSectionProps } from '../types';
	import { getMessageState } from './state.svelte';

	let { entityId }: PanelSectionProps = $props();

	const ms = getMessageState(entityId);

	function handlCompact() {
		ms.compactSignals();
	}

	function handleBulkDelete(signals: Signal[]) {
		ms.deleteSignals(signals.map((sig) => sig.entityId));
	}

	function handleDelete(signal: Signal) {
		ms.deleteSignal(signal.entityId);
	}
</script>

{#snippet section(msg: Message)}
	{#if msg.signals}
		<div class="flex flex-col gap-5 @5xl:gap-2 @5xl:flex-row">
			<div class="flex-1">
				<Table items={msg.signals} idKey="entityId">
					{#snippet bulkActions({ selectedCount, selectedItems, deselectAll })}
						<div class="flex justify-end gap-5">
							<IconButton onclick={handlCompact} color="secondary">
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

						<TableTitle>Size</TableTitle>

						<TableTitle>Start Position</TableTitle>
					{/snippet}

					{#snippet row(sig)}
						<TableField>{sig.name}</TableField>

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

	<pre>{JSON.stringify(msg, null, 2)}</pre>
{/snippet}

<section>
	<h3 class="pb-5">Signals</h3>

	{@render section(ms.entity)}
</section>
