<template>
	<div class="admin-add-question">
		<h1>問題列表</h1>
		<div class="table-scroll">
			<table>
				<thead>
					<tr>
						<th>ID</th>
						<th>問題</th>
						<th>選項</th>
					</tr>
				</thead>
				<tbody>
					<tr v-for="q in questions" :key="q.id">
						<td>{{ q.id }}</td>
						<td>{{ q.question }}</td>
						<td>{{ Array.isArray(q.options) ? q.options.join(', ') : '' }}</td>
					</tr>
				</tbody>
			</table>
		</div>

		<h2>新增問題</h2>
		<form @submit.prevent="submitQuestion">
			<label for="questionText">問題敘述</label>
			<textarea id="questionText" v-model="newQuestion" rows="2" required></textarea>

			<div class="options-list">
				<label>選項</label>
				<div v-for="(option, index) in options" :key="index" class="option-item">
					<input type="text" v-model="options[index]" :placeholder="'選項' + (index + 1)" required>
					<button type="button" @click="removeOption(index)" v-if="options.length > 1">刪除</button>
				</div>
				<button type="button" @click="addOption">新增選項</button>
			</div>

			<button type="submit">送出問題</button>
		</form>
	</div>
</template>

<script>
import http from "../../../router/axios";

export default {
	name: "AdminAddQuestion",
	data() {
		return {
			questions: [],
			newQuestion: "",
			options: [""]
		};
	},
	methods: {
		async loadQuestions() {
			await http.get("/question/")
				.then(response => {
					const data = response.data;
					if (data && data.data) {
						this.questions = data.data;
					}
				})
				.catch(err => console.error("載入問題錯誤:", err));
		},
		addOption() {
			this.options.push("");
		},
		removeOption(index) {
			this.options.splice(index, 1);
		},
		async submitQuestion() {
			const payload = {
				question: this.newQuestion,
				options: this.options.filter(o => o.trim() !== "")
			};

			await http.post("/question/", payload, {
				headers: {
					"Content-Type": "application/json"
				}
			})
			.then(response => {
				alert(response.data.message);
				this.loadQuestions();
				this.newQuestion = "";
				this.options = [""];
			})
			.catch(err => console.error("新增問題錯誤:", err));
		}
	},
	mounted() {
		this.loadQuestions();
	}
};
</script>

<style scoped>
.admin-add-question {
	padding: 20px;
	font-family: Arial, sans-serif;
}

.table-scroll {
	max-height: 400px;
	overflow-y: auto;
	margin-top: 20px;
	border: 1px solid #ccc;
}

table {
	width: 100%;
	border-collapse: collapse;
}

th,
td {
	border: 1px solid #ccc;
	padding: 8px;
	text-align: left;
}

form {
	margin-top: 30px;
	border: 1px solid #ccc;
	padding: 15px;
	max-width: 600px;
}

input,
textarea {
	width: 100%;
	padding: 10px;
	margin-bottom: 10px;
	box-sizing: border-box;
}

.options-list {
	margin-bottom: 10px;
}

.option-item {
	display: flex;
	align-items: center;
	margin-bottom: 5px;
}

.option-item input {
	flex-grow: 1;
}

.option-item button {
	margin-left: 10px;
}
</style>
