<script lang="ts">
	import type { Bus } from '$lib/api/canturin';
	import { UnderlinedButton } from '$lib/components/button';
	import { Table, TableTitle } from '$lib/components/table';
	import TableField from '$lib/components/table/table-field.svelte';
	import layoutStateSvelte from '$lib/state/layout-state.svelte';
	import { getHexNumber } from '$lib/utils';
	import type { PanelSectionProps } from '../types';
	import { getBusState } from './state.svelte';

	let { entityId }: PanelSectionProps = $props();

	const bs = getBusState(entityId);
</script>

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
					<UnderlinedButton
						label={node.name}
						onclick={() => layoutStateSvelte.openPanel('node', node.entityId)}
					/>
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
