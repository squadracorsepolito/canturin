<script lang="ts">
	import type { MultiplexerSignal } from '$lib/api/canturin';
	import { Attribute, AttributeGroup } from '$lib/components/attribute';
	import { Readonly } from '$lib/components/readonly';

	type Props = {
		entityId: string;
		signal: MultiplexerSignal;
	};

	let { entityId, signal }: Props = $props();
</script>

<AttributeGroup>
	<Attribute label="Group Count" desc="The number of groups">
		<Readonly>{signal.groupCount}</Readonly>
	</Attribute>

	<Attribute label="Group Size" desc="The size of each group in bits">
		<Readonly>{signal.groupSize}</Readonly>
	</Attribute>
</AttributeGroup>
