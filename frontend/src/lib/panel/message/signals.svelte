<script lang="ts">
	import { type Message, type BaseSignal, SignalKind } from '$lib/api/canturin';
	import { IconButton } from '$lib/components/button';
	import { SignalGrid } from '$lib/components/grid';
	import { SignalHoverPreview } from '$lib/components/hover-preview';
	import { AddIcon, CompactIcon, DeleteIcon } from '$lib/components/icon';
	import { Table, TableField, TableTitle } from '$lib/components/table';
	import type { PanelSectionProps } from '../types';
	import { getMessageState } from './state.svelte';
	import { AddSignalModal } from '$lib/components/modal';
	import { openPanel } from '../panel-stack-state.svelte';
	import { SignalKindBadge } from '$lib/components/badge';

	let { entityId }: PanelSectionProps = $props();

	const ms = getMessageState(entityId);

	function handleAdd(signalKind: SignalKind) {
		ms.addSignal(signalKind);
	}

	function handleCompact() {
		ms.compactSignals();
	}

	function handleReorder(sigEntId: string, from: number, to: number) {
		ms.reorderSignal(sigEntId, from, to);
	}

	function handleBulkDelete(signals: BaseSignal[]) {
		ms.deleteSignals(signals.map((sig) => sig.entityId));
	}

	function handleDelete(signal: BaseSignal) {
		ms.deleteSignal(signal.entityId);
	}
</script>

{#snippet section(msg: Message)}
	{#if msg.signals}
		<div class="flex flex-col gap-5 @5xl:gap-2 @5xl:flex-row">
			<div class="flex-1">
				<Table items={msg.signals} idKey="entityId" reorder={handleReorder}>
					{#snippet bulkActions({ selectedCount, selectedItems, deselectAll })}
						<div class="flex justify-end gap-5">
							<IconButton onclick={handleCompact} themeColor="secondary">
								<CompactIcon />
							</IconButton>

							<AddSignalModal onsubmit={handleAdd}>
								{#snippet trigger({ getProps })}
									<IconButton
										label="Add Signal"
										themeColor="primary"
										disabled={msg.maxAvailableSpace === 0}
										{...getProps()}
									>
										<AddIcon />
									</IconButton>
								{/snippet}
							</AddSignalModal>

							<IconButton
								onclick={() => {
									handleBulkDelete(selectedItems);
									deselectAll();
								}}
								label={`Delete Signals ${selectedCount > 0 ? ` (${selectedCount})` : ''}`}
								disabled={selectedCount === 0}
								themeColor="error"
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
							<SignalHoverPreview signal={sig} />
						</TableField>

						<TableField>
							<SignalKindBadge signalKind={sig.kind} />
						</TableField>

						<TableField>{sig.size}</TableField>

						<TableField>{sig.startPos}</TableField>
					{/snippet}

					{#snippet rowActions(signal)}
						<IconButton onclick={() => handleDelete(signal)} themeColor="error">
							<DeleteIcon />
						</IconButton>
					{/snippet}
				</Table>
			</div>

			<div class="flex-1">
				<SignalGrid
					signals={msg.signals}
					height={msg.sizeByte}
					onclick={(entityId, name) => openPanel('signal', entityId, name)}
				/>
			</div>
		</div>
	{/if}
{/snippet}

<section>
	<h3 class="pb-5">Signals</h3>

	{@render section(ms.entity)}
</section>
