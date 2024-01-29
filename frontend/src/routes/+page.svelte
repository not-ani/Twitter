<script lang="ts">
	import { onMount } from 'svelte';
	import type { InnerTweets } from '../types';
	let data: InnerTweets[] = [];

	// your script goes here
	onMount(async () => {
		const response = await fetch(`http://localhost:3000/api/tweet/all/`);
		data = (await response.json()) as InnerTweets[];
	});
</script>

<div>
	<h1>Hello World</h1>
	{#each data as tweet}
		<a href={`/tweet/${tweet.id}`}>
			<p>{tweet.content}</p>
		</a>
	{/each}
</div>
