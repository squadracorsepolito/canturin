<script lang="ts">
	import type { BaseSignal, MultiplexerSignalGroup } from '$lib/api/canturin';
	import { SignalKindBadge } from '$lib/components/badge';
	import { IconButton } from '$lib/components/button';
	import { SignalGrid } from '$lib/components/grid';
	import { SignalHoverPreview } from '$lib/components/hover-preview';
	import { DeleteIcon } from '$lib/components/icon';
	import { Pagination } from '$lib/components/pagination';
	import { Table, TableField, TableTitle } from '$lib/components/table';
	import { getSignalState } from './state.svelte';

	type Props = {
		entityId: string;
	};

	let { entityId }: Props = $props();

	const ss = getSignalState(entityId);

	function handleBulkDelete(groupId: number, signals: BaseSignal[]) {
		ss.deleteMultiplexedSignals(
			groupId,
			signals.map((s) => s.entityId)
		);
	}

	function handleDelete(groupId: number, signal: BaseSignal) {
		ss.deleteMultiplexedSignal(groupId, signal.entityId);
	}

	let groupId = $state(0);

	let emptyGroupSet = $derived.by(() => {
		if (!ss.entity.multiplexer.emptyGroups) return new Set<number>();

		return new Set(ss.entity.multiplexer.emptyGroups);
	});
</script>

{#snippet singleGroup(group: MultiplexerSignalGroup, groupSize: number)}
	{#if group.signals}
		<div class="flex flex-col gap-5 @5xl:gap-2 @5xl:flex-row">
			<div class="flex-1">
				<Table items={group.signals} idKey="entityId">
					{#snippet bulkActions({ selectedCount, selectedItems, deselectAll })}
						<div class="flex justify-end gap-5">
							<IconButton
								onclick={() => {
									handleBulkDelete(group.id, selectedItems);
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

					{#snippet row(muxSig)}
						<TableField>
							<SignalHoverPreview signal={muxSig} />
						</TableField>

						<TableField>
							<SignalKindBadge signalKind={muxSig.kind} />
						</TableField>

						<TableField>
							{muxSig.size}
						</TableField>

						<TableField>
							{muxSig.startPos}
						</TableField>
					{/snippet}

					{#snippet rowActions(muxSig)}
						<IconButton onclick={() => handleDelete(group.id, muxSig)} themeColor="error">
							<DeleteIcon />
						</IconButton>
					{/snippet}
				</Table>
			</div>

			<div class="flex-1">
				<SignalGrid width={8} height={groupSize / 8} signals={group.signals} />
			</div>
		</div>
	{/if}
{/snippet}

<section>
	<h3 class="pb-5">Groups</h3>

	<div class="flex justify-center pb-5">
		<Pagination bind:selectedPage={groupId} totPages={ss.entity.multiplexer.groupCount} />
	</div>

	{#if ss.entity.multiplexer.groups}
		{@const groupSize = ss.entity.multiplexer.groupSize}

		{#if emptyGroupSet.has(groupId)}
			{@render singleGroup({ id: groupId, signals: [] }, groupSize)}
		{:else}
			{@const group = ss.entity.multiplexer.groups[groupId]}

			{@render singleGroup(group, groupSize)}
		{/if}
	{/if}
</section>

<!-- 
<section>
	<h3 class="pb-5">Multiplexed Signals</h3>

	{#if muxSignals}
		<Table items={muxSignals} idKey="entityId">
			{#snippet header()}
				<TableTitle>Group ID</TableTitle>

				<TableTitle>Name</TableTitle>

				<TableTitle>Kind</TableTitle>

				<TableTitle>Size</TableTitle>

				<TableTitle>Start Position</TableTitle>
			{/snippet}

			{#snippet row(muxSig)}
				<TableField>{muxSig.groupId}</TableField>

				<TableField>
					<SignalHoverPreview signal={muxSig} />
				</TableField>

				<TableField>
					<SignalKindBadge signalKind={muxSig.kind} />
				</TableField>

				<TableField>
					{muxSig.size}
				</TableField>

				<TableField>
					{muxSig.startPos}
				</TableField>
			{/snippet}

			{#snippet rowActions(muxSig)}
				<IconButton onclick={() => console.log('delete')} themeColor="error">
					<DeleteIcon />
				</IconButton>
			{/snippet}
		</Table>
	{/if}
</section> -->
