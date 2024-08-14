<script lang="ts">
	import { SignalUnitService, type SignalUnit } from '$lib/api/canturin';
	import Table from '$lib/components/table/table.svelte';
	import Panel from './panel.svelte';

	type Props = {
		entityId: string;
	};

	let { entityId }: Props = $props();

	let sigUnitPromise = $derived.by(async () => {
		try {
			const tmpSigUnit = await SignalUnitService.GetOpen(entityId);
			return tmpSigUnit;
		} catch (error) {
			console.error(error);
		}
	});
</script>

{#snippet sigUnitPanel(sigUnit: SignalUnit)}
	<section>
		<h3>{sigUnit.name}</h3>
		<p>{sigUnit.desc}</p>
	</section>

	{#if sigUnit.references}
		<section>
			<h4>References</h4>

			<Table items={sigUnit.references}>
				{#snippet header()}
					<th>Message</th>
					<th>Signals</th>
				{/snippet}

				{#snippet row(ref)}
					<td>{ref.parentMessage.name}</td>
					<td>{ref.name}</td>
				{/snippet}
			</Table>
		</section>
	{/if}
{/snippet}

<Panel>
	{#await sigUnitPromise then sigUnit}
		{#if sigUnit}
			{@render sigUnitPanel(sigUnit)}
		{/if}
	{/await}
</Panel>
