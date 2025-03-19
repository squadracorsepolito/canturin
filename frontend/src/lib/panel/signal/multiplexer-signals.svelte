<script lang="ts">
	import type { MultiplexerSignal } from '$lib/api/canturin';
	import { SignalHoverPreview } from '$lib/components/hover-preview';
	import { Table, TableField, TableTitle } from '$lib/components/table';

	type Props = {
		entityId: string;
		signal: MultiplexerSignal;
	};

	let { entityId, signal }: Props = $props();

	let muxSignals = $derived.by(() => {
		if (!signal.groups) return;

		const filteredGroups = signal.groups.filter((g) => {
			if (!g.signals) return false;

			return g.signals.length > 0;
		});

		const res = [];
		for (const group of filteredGroups) {
			const groupId = group.id;
			for (const muxSig of group.signals || []) {
				res.push({
					groupId,
					...muxSig
				});
			}
		}
		return res;
	});
</script>

<section>
	<h3 class="pb-5">Multiplexed Signals</h3>

	{#if muxSignals}
		<Table items={muxSignals} idKey="entityId">
			{#snippet header()}
				<TableTitle>Group ID</TableTitle>

				<TableTitle>Name</TableTitle>
			{/snippet}

			{#snippet row(muxSig)}
				<TableField>{muxSig.groupId}</TableField>

				<TableField>
					<SignalHoverPreview signal={muxSig} />
				</TableField>
			{/snippet}
		</Table>
	{/if}
</section>
