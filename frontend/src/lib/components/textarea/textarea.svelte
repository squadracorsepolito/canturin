<script lang="ts">
	type Props = {
		value: string;
	};

	let { value = $bindable() }: Props = $props();
</script>

<textarea bind:value class="textarea textarea-bordered w-full"></textarea>
