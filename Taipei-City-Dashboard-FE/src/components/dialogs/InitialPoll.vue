<!-- Developed by Taipei Urban Intelligence Center 2023-2024-->

<script setup>
import { ref, onMounted } from "vue";
import { useDialogStore } from "../../store/dialogStore";
import { useAuthStore } from "../../store/authStore";
import DialogContainer from "./DialogContainer.vue";
import http from "../../router/axios"; // Make sure this is your axios instance

const pollID = ref(5); // Default poll ID, can be changed if needed
const pollTitle = ref("請表態您的立場");
const pollOptions = ref(["Option 1", "Option 2"]); // default fallback

onMounted(async () => {
	try {
		const response = await http.get(`/question/${pollID.value}`);
		if (response.data) {
			pollTitle.value = response.data.data.question || pollTitle.value;
			pollOptions.value = response.data.data.options || pollOptions.value;
			console.log("Poll data fetched successfully:", response.data.data);
		}
	} catch (error) {
		console.error("Failed to fetch poll data:", error);
	}
});
const dialogStore = useDialogStore();
const authStore = useAuthStore();

const dontShowAgain = ref(false);

function handleClose() {
	dialogStore.hideAllDialogs();
}

function handleSubmit() {
	if (dontShowAgain.value) {
		localStorage.setItem("initialWarning", "shown");
	}
	handleClose();
}

async function handleVote(option) {
	// handle vote logic here
	console.log("Voted for:", option);
	const optionIndex = pollOptions.value.indexOf(option);
	console.log("Voted for option index:", optionIndex);
	await http.post(`/question/${pollID.value}/vote`, {
		choice: optionIndex
	});
}
</script>

<template>
	<DialogContainer dialog="initialWarning" @on-close="handleClose">
		<div class="login">
			<div class="login-logo">
				<div class="login-logo-image">
					<img
						src="../../assets/images/TUIC.svg"
						alt="tuic logo"
						@click.shift="handleSwitchMode"
					/>
				</div>
				<div>
					<!-- <h1>{{ VITE_APP_TITLE }}</h1> -->
					<h2>Taipei City Dashboard Poll</h2>
					<h1>雙北城市儀表板</h1>
				</div>
			</div>
			<!-- Added poll section -->
			<div class="poll-section">
				<h3 class="poll-title">{{ pollTitle }}</h3>
				<div class="poll-options">
					<button
						v-for="(option, idx) in pollOptions"
						:key="idx"
						@click="handleVote(option)"
					>
						{{ option }}
					</button>
				</div>
			</div>

			<p>
				<a
					href="https://tuic.gov.taipei/zh/works/dashboard"
					target="_blank"
					>臺北城市儀表板</a
				>的<a href="https://tuic.gov.taipei/zh/privacy" target="_blank"
					>隱私權政策</a
				>
			</p>
			<p
				:style="{
					color: '#302C2E',
					cursor: 'default',
					userSelect: 'none',
				}"
			>
				TUIC Igor Ann Iima Chu Jack 2023-2024
			</p>
			<p>《讓臺北城市儀表板成為您的儀表板》</p>
		</div>
	</DialogContainer>
</template>

<style scoped lang="scss">
.poll-section {
	margin: 20px 0;
	text-align: center;
	width: 100%;

	.poll-title {
		font-size: 1.5rem;
		color: #fff;
		margin-bottom: 28px;
	}

	.poll-options {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 20px; // Add more space between buttons
		margin: 0 auto;
		max-width: 400px;

		button {
			all: unset;
			display: block;
			width: 90%;
			max-width: 340px;
			padding: 12px 0;
			margin: 0 auto;
			font-size: 1.1rem;
			color: #fff;
			background: #03b2c3;
			border-radius: 16px;
			text-align: center;
			cursor: pointer;
			transition: background 0.2s, color 0.2s, box-shadow 0.2s;
			box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);

			&:hover {
				background: var(--color-highlight);
				color: #fff;
			}
		}
	}
}

.login {
	width: 300px;

	p {
		text-align: center;
		color: var(--color-complement-text);
	}

	p:last-child {
		background: linear-gradient(
			75deg,
			var(--color-complement-text),
			var(--color-highlight) 70%,
			var(--color-complement-text)
		);
		background-clip: text;
		-webkit-background-clip: text;
		-webkit-text-fill-color: transparent;
		animation: title-gradient 10s linear none infinite;
	}

	a {
		color: var(--color-highlight);
	}

	label {
		margin-bottom: 4px;
		color: var(--color-complement-text);
		font-size: var(--font-s);
		align-self: flex-start;
	}

	input {
		margin-bottom: 8px;
		width: calc(100% - 14px);
	}

	&-logo {
		display: flex;
		justify-content: center;

		h1 {
			font-weight: 500;
		}

		h2 {
			font-size: var(--font-s);
			font-weight: 400;
		}

		&-image {
			width: 22.94px;
			height: 45px;
			margin: 0 10px 0 0;

			img {
				height: 45px;
				filter: invert(1);
			}
		}
	}

	&-form {
		width: 100%;
		height: 200px;
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
	}
}

@keyframes title-gradient {
	0% {
		background-position: 0;
	}

	50% {
		background-position: 600px;
	}

	100% {
		background-position: 0;
	}
}
</style>
