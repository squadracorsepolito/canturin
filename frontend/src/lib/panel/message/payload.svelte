<script lang="ts">
	import { colorByName } from '$lib/actions/color-name.svelte';
	import type { Message } from '$lib/api/canturin';
	import { ForwardIcon } from '$lib/components/icon';
	import BackwardIcon from '$lib/components/icon/backward-icon.svelte';
	import { Table, TableField, TableTitle } from '$lib/components/table';
	import type { PanelSectionProps } from '../types';
	import { getMessageState, type GridItem, gridWidth } from './state.svelte';

	let { entityId }: PanelSectionProps = $props();

	const ms = getMessageState(entityId);

	function getPlaceholder(idx: number) {
		let nRow = Math.ceil(idx / gridWidth);
		if (idx % gridWidth === 0) {
			nRow++;
		}
		return gridWidth - idx + 2 * gridWidth * (nRow - 1) - 1;
	}
</script>

{#snippet gridItem({ colStart, colEnd, rowStart, rowEnd, id, name, continues, follows }: GridItem)}
	<div
		use:colorByName={{ name }}
		onpointerenter={() => (ms.hoveredID = id)}
		onpointerleave={() => (ms.hoveredID = '')}
		class="rounded-btn text-xs cursor-pointer hover:opacity-50 transition-all {id ===
			ms.hoveredID && 'opacity-50'}"
		style:grid-column-start={colStart}
		style:grid-column-end={colEnd}
		style:grid-row-start={rowStart}
		style:grid-row-end={rowEnd}
	>
		<div class="flex p-3 h-full w-full">
			{#if ms.hoveredID !== id}
				{#if continues}
					<span class="self-end"><ForwardIcon /></span>
				{/if}

				<span class="flex-1 self-center text-center truncate text-ellipsis">
					{name}
				</span>

				{#if follows}
					<span class="self-start"> <BackwardIcon /></span>
				{/if}
			{/if}
		</div>
	</div>
{/snippet}

{#snippet grid(sizeByte: number, items: GridItem[])}
	<div class="flex gap-2">
		<div class="flex flex-col py-2 gap-2 justify-around">
			{#each { length: sizeByte } as _, idx}
				<div class="font-semibold pr-5">
					{idx + 1}
				</div>
			{/each}
		</div>

		<div class="flex-1 relative">
			<div class="absolute h-full w-full block rounded-box bg-base-200"></div>

			<div
				class="absolute w-full h-full grid gap-2 p-2"
				style:grid-template-columns="repeat({gridWidth}, minmax(0, 1fr))"
				style:grid-template-rows="repeat({sizeByte}, minmax(0, 1fr))"
			>
				{#each { length: gridWidth * sizeByte } as _, idx}
					<div class="flex items-center rounded-box justify-center bg-base-100">
						<span class="text-base-content font-semibold">
							{getPlaceholder(idx)}
						</span>
					</div>
				{/each}
			</div>

			<div
				class="relative grid gap-2 p-2"
				style:grid-template-columns="repeat({gridWidth}, minmax(0, 1fr))"
				style:grid-template-rows="repeat({sizeByte}, minmax(0, 1fr))"
				style:aspect-ratio={gridWidth / sizeByte}
			>
				{#each items as sigItem}
					{@render gridItem(sigItem)}
				{/each}
			</div>
		</div>
	</div>
{/snippet}

{#snippet section(msg: Message)}
	{#if msg.signals}
		<div class="flex flex-col gap-5 @5xl:gap-2 @5xl:flex-row">
			<div class="flex-1">
				<Table items={msg.signals} idKey="entityId">
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
				</Table>
			</div>

			<div class="flex-1">
				{@render grid(msg.sizeByte, ms.gridItems)}
			</div>
		</div>
	{/if}

	<pre>{JSON.stringify(msg, null, 2)}</pre>
{/snippet}

<section>
	<h3 class="pb-5">Payload</h3>

	{@render section(ms.entity)}
</section>
