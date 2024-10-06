<script lang="ts">
	type Props = {
		title: string;
		value: string | number;
		desc?: string;
		isEditing?: boolean;
	};

	let { title, value, desc, isEditing }: Props = $props();
</script>

<div
	class="flex flex-col items-center p-3 rounded-t-box {isEditing
		? 'bg-primary text-primary-content'
		: 'hover:bg-base-200 rounded-b-box'}"
>
	<div class="opacity-90">{title}</div>

	<div class="text-3xl font-extrabold">{value}</div>

	{#if desc}
		<div class="opacity-90 text-xs">{desc}</div>
	{/if}
</div>
