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
<style scoped lang="scss">
.admin-add-question {
  padding: 24px;
  background-color: #000;
  color: #fff;
  width: 100%;
  box-sizing: border-box;

  h1 {
    font-size: 24px;
    font-weight: bold;
    margin-bottom: 16px;
  }

  h2 {
    font-size: 20px;
    font-weight: bold;
    margin-top: 32px;
    margin-bottom: 16px;
  }

  .table-scroll {
    max-height: 400px;
    overflow-y: auto;
    margin-top: 20px;
    border: 1px solid #333;
    border-radius: 4px;
    width: 100%;

    table {
      width: 100%;
      border-collapse: collapse;

      thead {
        background-color: #333;
        th {
          padding: 12px;
          border-bottom: 1px solid #333;
          text-align: left;
          color: #fff;
        }
      }

      tbody {
        tr {
          &:hover {
            background-color: #444;
          }
          td {
            padding: 12px;
            border-bottom: 1px solid #333;
            text-align: left;
            color: #fff;
          }
        }
      }
    }
  }

  form {
    margin-top: 30px;

    label {
      font-weight: bold;
      display: block;
      margin-bottom: 8px;
      font-size: 14px;
      color: #fff;
    }

    textarea,
    input[type="text"] {
      width: 100%;
      padding: 10px;
      margin-bottom: 16px;
      border: 1px solid #555;
      border-radius: 4px;
      font-size: 15px;
      box-sizing: border-box;
      background-color: #222;
      color: #fff;
    }

    .options-list {
      .option-item {
        display: flex;
        align-items: center;
        margin-bottom: 10px;

        input {
          flex: 1;
        }

        button {
          margin-left: 10px;
          padding: 6px 10px;
          border: none;
          border-radius: 4px;
          background-color: #e57373;
          color: #fff;
          cursor: pointer;

          &:hover {
            background-color: #d32f2f;
          }
        }
      }

      > button {
        padding: 10px 16px;
        border: none;
        border-radius: 4px;
        background-color: #42a5f5;
        color: #fff;
        cursor: pointer;
        font-size: 15px;
        margin-top: 10px;
        width: 100%;

        &:hover {
          background-color: #1e88e5;
        }
      }
    }

    > button[type="submit"] {
      padding: 10px 16px;
      border: none;
      border-radius: 4px;
      background-color: #42a5f5;
      color: #fff;
      cursor: pointer;
      font-size: 15px;
      margin-top: 10px;
      width: 100%;

      &:hover {
        background-color: #1e88e5;
      }
    }
  }

  @media (max-width: 520px) {
    padding: 16px;
  }
}
</style>