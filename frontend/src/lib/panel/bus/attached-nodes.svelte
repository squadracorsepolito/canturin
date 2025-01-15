<script lang="ts">
	import type { BaseEntity, Bus } from '$lib/api/canturin';
	import { LinkButton } from '$lib/components/button';
	import { HoverPreview } from '$lib/components/hover-preview';
	import { Table, TableTitle } from '$lib/components/table';
	import TableField from '$lib/components/table/table-field.svelte';
	import layoutStateSvelte from '$lib/state/layout-state.svelte';
	import { getHexNumber } from '$lib/utils';
	import type { PanelSectionProps } from '../types';
	import { getBusState } from './state.svelte';

	let { entityId }: PanelSectionProps = $props();

	const bs = getBusState(entityId);
</script>

{#snippet nodePreview(node: BaseEntity)}
	<div class="font-medium text-sm pr-1">{node.name}</div>

	{#if node.desc}
		<div class="text-xs text-dimmed pt-1">{node.desc}</div>
	{/if}
{/snippet}

{#snippet section(bus: Bus)}
	{#if bus.attachedNodes}
		<Table items={bus.attachedNodes} idKey="entityId">
			{#snippet header()}
				<TableTitle>Name</TableTitle>

				<TableTitle>ID</TableTitle>

				<TableTitle>Hex ID</TableTitle>

				<TableTitle>Interface Number</TableTitle>
			{/snippet}

			{#snippet row(node)}
				<TableField>
					<HoverPreview placement="right">
						{#snippet trigger()}
							<LinkButton
								label={node.name}
								onclick={() => layoutStateSvelte.openPanel('node', node.entityId)}
							/>
						{/snippet}

						{#snippet content()}
							{@render nodePreview(node)}
						{/snippet}
					</HoverPreview>
				</TableField>

				<TableField>{node.id}</TableField>

				<TableField>{getHexNumber(node.id)}</TableField>

				<TableField>{node.interfaceNumber}</TableField>
			{/snippet}
		</Table>
	{/if}
{/snippet}

<section>
	<h3 class="pb-5">Attached Nodes</h3>

	{@render section(bs.entity)}
</section>
