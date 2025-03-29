<script lang="ts">
	import type { Ref } from '$lib/api/canturin';
	import ReferenceTree1 from '$lib/components/tree/reference-tree1.svelte';

	type Props = {
		reference: Ref;
	};

	let { reference }: Props = $props();
</script>

<section>
	<h3 class="pb-5">References</h3>

	<ReferenceTree1 root={reference} />
</section>
