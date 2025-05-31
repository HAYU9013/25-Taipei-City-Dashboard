<!-- Developed by Taipei Urban Intelligence Center 2023-2024 -->

<script setup>
import { ref, onMounted } from "vue";
import { useDialogStore } from "../../store/dialogStore";
import { useAuthStore } from "../../store/authStore";
import DialogContainer from "./DialogContainer.vue";
import http from "../../router/axios"; // Make sure this is your axios instance
import DashboardComponent from "../../dashboardComponent/DashboardComponent.vue";

const pollID = ref(5); // Default poll ID, can be changed if needed
const pollTitle = ref("請表態您的立場");
const pollOptions = ref(["Option 1", "Option 2"]); // default fallback

onMounted(async () => {
	try {
		const response = await http.get("/question/");
		if (response.data && Array.isArray(response.data.data) && response.data.data.length) {
			const questions = response.data.data;
			const latest = questions.reduce((prev, curr) => (curr.id > prev.id ? curr : prev), questions[0]);
			pollID.value = latest.id;
			console.log("Latest poll id set to:", pollID.value);
		}
	} catch (error) {
		console.error("Failed to fetch questions:", error);
	}
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

// 用來存放從後端拿到的完整「config」
const dashboardConfig = ref(null);
// 用來存放從後端拿到的「外層陣列 of 物件」，每個物件底下有 .data
const dashboardChartData = ref([]);

// 是否顯示 DashboardComponent 與其設定資料
const showDashboard = ref(false);


const dialogStore = useDialogStore();
const authStore = useAuthStore();

onMounted(async () => {
  // 先取得 Poll 的標題與選項
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

  // 再取得儀表板的 Component 設定 (config)
  try {
    const configRes = await http.get(`/component/4?city=metrotaipei`);
    if (configRes.data && configRes.data.data) {
      dashboardConfig.value = { ...configRes.data.data };
      console.log("▶️ dashboardConfig:", dashboardConfig.value);
    }
  } catch (error) {
    console.error("取得 component 設定失敗：", error);
  }
});

async function fetchChartData() {
  try {
    const chartRes = await http.get(`/component/4/chart?city=metrotaipei`);
    console.log("chartRes =", chartRes);
    if (
      chartRes.data &&
      Array.isArray(chartRes.data.data) &&
      chartRes.data.data.length > 0
    ) {
      // 將外層陣列指派給 dashboardChartData
      dashboardChartData.value = chartRes.data.data;
      console.log("▶️ dashboardChartData:", dashboardChartData.value);
    } else {
      dashboardChartData.value = [];
    }
  } catch (error) {
    console.error("取得圖表資料失敗：", error);
    dashboardChartData.value = [];
  }
}

function handleClose() {
  dialogStore.hideAllDialogs();
}

async function handleVote(option) {
  console.log("Voted for:", option);
  const optionIndex = pollOptions.value.indexOf(option) + 1;
  try {
    await http.post(`/question/${pollID.value}/vote`, {
      choice: optionIndex,
    });
  } catch (error) {
    console.error("投票失敗：", error);
  }
  // 投票成功後顯示儀表板，並且去撈圖表資料
  showDashboard.value = true;
  await fetchChartData();
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
          <h2>Taipei City Dashboard Poll</h2>
          <h1>雙北城市儀表板</h1>
        </div>
      </div>

      <!-- 如果還沒投票，顯示投票選項 -->
      <div v-if="!showDashboard">
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
      </div>

      <!-- 投票完成後，顯示 DashboardComponent -->
      <div v-else>
        <div v-if="dashboardConfig && dashboardChartData.length > 0">
          <DashboardComponent
            :config="{
              ...dashboardConfig,
              // 注意：chart_data 必須是「外層陣列 of 物件」
              chart_data: dashboardChartData
            }"
            mode="default"
          />
        </div>
        <div v-else>
          <p>資料讀取中，請稍候…</p>
        </div>
        <p>按空白處退出</p>
      </div>

      <p>
        <a href="https://tuic.gov.taipei/zh/works/dashboard" target="_blank"
          >臺北城市儀表板</a
        >的
        <a href="https://tuic.gov.taipei/zh/privacy" target="_blank"
          >隱私權政策</a
        >
      </p>
      <p
        :style="{
          color: '#302C2E',
          cursor: 'default',
          userSelect: 'none'
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
    gap: 20px; /* Add more space between buttons */
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
