<!-- Developed By Taipei Urban Intelligence Center 2023-2024 -->
<!-- 
Lead Developer:  Igor Ho (Full Stack Engineer)
Data Pipelines:  Iima Yu (Data Scientist)
Design and UX: Roy Lin (Fmr. Consultant), Chu Chen (Researcher)
Systems: Ann Shih (Systems Engineer)
Testing: Jack Huang (Data Scientist), Ian Huang (Data Analysis Intern) 
-->
<!-- Department of Information Technology, Taipei City Government -->

<script setup>
// import * as THREE from 'three'
// import { GLTFLoader } from 'three/examples/jsm/loaders/GLTFLoader.js'
// import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls.js'
// import { nextTick } from 'vue'

// const threeContainer = ref(null)
// let renderer, scene, camera, model, controls, animationId

// const showThree = ref(false)

// onMounted(() => {
//   showThree.value = localStorage.getItem("show3DTest") === "true"

//   if (showThree.value) {
//     nextTick(() => {
//       if (!threeContainer.value) {
//         console.error("threeContainer not ready")
//         return
//       }

//       // Scene & Camera
//       scene = new THREE.Scene()
//       scene.background = new THREE.Color(0x111111) // dark background
//       camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 1000)
//       camera.position.set(0, 2, 5)

//       // Renderer
//       renderer = new THREE.WebGLRenderer({ antialias: true, alpha: true })
//       renderer.setSize(window.innerWidth, window.innerHeight)
//       renderer.shadowMap.enabled = true
//       threeContainer.value.appendChild(renderer.domElement)

//       // Controls
//       controls = new OrbitControls(camera, renderer.domElement)
//       controls.enableDamping = true
//       controls.dampingFactor = 0.05

//       // Lights
//       const ambientLight = new THREE.AmbientLight(0xffffff, 0.5)
//       scene.add(ambientLight)

//       const spotLight = new THREE.SpotLight(0xffffff, 2)
//       spotLight.position.set(5, 10, 5)
//       spotLight.castShadow = true
//       scene.add(spotLight)

//       // Shadow Plane
//       const planeGeometry = new THREE.PlaneGeometry(500, 500)
//       const planeMaterial = new THREE.ShadowMaterial({ opacity: 0.3 })
//       const shadowPlane = new THREE.Mesh(planeGeometry, planeMaterial)
//       shadowPlane.rotation.x = -Math.PI / 2
//       shadowPlane.position.y = -1
//       shadowPlane.receiveShadow = true
//       scene.add(shadowPlane)

//       // Load Model
//       const loader = new GLTFLoader()
//       loader.load('/models/Codefest_Penguin_0531155635_texture.glb', (gltf) => {
//         model = gltf.scene
//         model.scale.set(1, 1, 1)
//         model.traverse(obj => {
//           if (obj.isMesh) {
//             obj.castShadow = true
//             obj.receiveShadow = true
//           }
//         })
//         scene.add(model)
//         animate()
//       })
//     })
//   }
// })

// function animate() {
//   animationId = requestAnimationFrame(animate)
//   if (model) model.rotation.y += 0.025
//   controls.update()
//   renderer.render(scene, camera)
// }

// onBeforeUnmount(() => {
//   cancelAnimationFrame(animationId)
//   renderer?.dispose?.()
// })



import ThreePenguinShow from './dashboardComponent/components/ThreePenguinShow.vue'

import { onBeforeMount, onMounted, onBeforeUnmount, ref, computed, watch } from "vue";
import { useRoute } from "vue-router"
import { useAuthStore } from "./store/authStore";
import { useDialogStore } from "./store/dialogStore";
import { useContentStore } from "./store/contentStore";
import { useMapStore } from "./store/mapStore";

import NavBar from "./components/utilities/bars/NavBar.vue";
import SideBar from "./components/utilities/bars/SideBar.vue";
import AdminSideBar from "./components/utilities/bars/AdminSideBar.vue";
import SettingsBar from "./components/utilities/bars/SettingsBar.vue";
import NotificationBar from "./components/dialogs/NotificationBar.vue";
import InitialWarning from "./components/dialogs/InitialWarning.vue";
import ComponentSideBar from "./components/utilities/bars/ComponentSideBar.vue";
import LogIn from "./components/dialogs/LogIn.vue";
import InitialPoll from "./components/dialogs/InitialPoll.vue";

const authStore = useAuthStore();
const dialogStore = useDialogStore();
const contentStore = useContentStore();
const timeToUpdate = ref(600);

const mapStore = useMapStore();
const route = useRoute();
const updateBoards = import.meta.env.VITE_PERSONAL_BOARD_UPDATE?.split(',') || [];
const boardIndex = ref(null);
const board =ref(null);
const frequency = ref(600);
const isMappedToUpdateBoards =ref(false);

const updateBoardsMap = computed(()=>{
	let needUpdateBoards = []
	updateBoards.map((board) => {
		const id = board.split(':')[0];
		const updateSeconds = board.split(':')[1];
		needUpdateBoards.push({ id, frequency: updateSeconds });
	})
	return needUpdateBoards
})

const formattedTimeToUpdate = computed(() => {
	const minutes = Math.floor(timeToUpdate.value / 60);
	const seconds = timeToUpdate.value % 60;
	return `${minutes}:${seconds < 10 ? "0" : ""}${seconds}`;
});

function reloadChartData() {
	if (!["dashboard", "mapview"].includes(authStore.currentPath)) return;
	contentStore.setCurrentDashboardAllChartData();
	timeToUpdate.value = frequency.value;

	if (isMappedToUpdateBoards.value) {
		reloadMapData()
	}
}
function updateTimeToUpdate() {
	if (!["dashboard", "mapview"].includes(authStore.currentPath)) return;
	if (timeToUpdate.value <= 0) {
		timeToUpdate.value = 0;
		reloadChartData();
		return;
	}
	timeToUpdate.value -= 5;
}

function reloadMapData() {
	if (!["mapview"].includes(authStore.currentPath)) return;
	mapStore.currentVisibleLayers.forEach((layerName) => {
		mapStore.map.removeLayer(layerName);
		if (mapStore.map.getSource(`${layerName}-source`)) {
			mapStore.map.removeSource(`${layerName}-source`);
		}
		const layerConfig = mapStore.mapConfigs[layerName];

		// 檢查 source
		if (layerConfig.source === "geojson") {
			// 如果 source 是 "geojson"，則使用 fetchLocalGeoJson
			mapStore.fetchLocalGeoJson(layerConfig);
		} else if (layerConfig.source === "raster") {
			// 如果 source 是 "raster"，則使用 addRasterSource
			mapStore.addRasterSource(layerConfig);
		}
	});
}

watch(() => route.query, (query) => {
	boardIndex.value = query.index;
	board.value = updateBoardsMap.value.find(board =>{
		return board.id === boardIndex.value
	})
	frequency.value = board.value ? board.value.frequency : 600;
	isMappedToUpdateBoards.value = updateBoardsMap.value.some(board =>{
		return board.id === query.index
	});
	timeToUpdate.value = frequency.value;

}), { immediate: true };

onBeforeMount(() => {
	authStore.initialChecks();

	let vh = window.innerHeight * 0.01;
	document.documentElement.style.setProperty("--vh", `${vh}px`);

	window.addEventListener("resize", () => {
		let vh = window.innerHeight * 0.01;
		document.documentElement.style.setProperty("--vh", `${vh}px`);
	});
	// contentStore.wsConnect();
});
onMounted(() => {
	const showInitialWarning = localStorage.getItem("initialWarning");
	
	if (!showInitialWarning && !window.location.pathname.includes("embed")) {
		dialogStore.showDialog("initialWarning");
	}

	setInterval(reloadChartData, 1000 * frequency.value);
	setInterval(updateTimeToUpdate, 1000 * 5);
});
onBeforeUnmount(() => {
	clearInterval(reloadChartData);
	clearInterval(updateTimeToUpdate);
	// contentStore.wsDisconnect();
});
</script>

<template>
  <div class="app-container">
    <NotificationBar />
    <NavBar v-if="authStore.currentPath !== 'embed'" />
    <!-- /mapview, /dashboard layouts -->
    <div
      v-if="
        authStore.currentPath === 'mapview' ||
          authStore.currentPath === 'dashboard'
      "
      class="app-content"
    >
      <SideBar />
      <div class="app-content-main">
        <SettingsBar />
        <RouterView />
      </div>
    </div>
    <!-- /admin layouts -->
    <div
      v-else-if="authStore.currentPath === 'admin'"
      class="app-content"
    >
      <AdminSideBar />
      <div class="app-content-main">
        <RouterView />
      </div>
    </div>
    <!-- /component, /component/:index layouts -->
    <div
      v-else-if="authStore.currentPath.includes('component')"
      class="app-content"
    >
      <ComponentSideBar />
      <div class="app-content-main">
        <RouterView />
      </div>
    </div>
    <div v-else>
      <router-view />
    </div>
    <!-- <InitialWarning /> -->
	<InitialPoll />
    <LogIn />
	<!-- ⬇️ Three.js 測試區，只在 test=3d 出現 -->
	<!-- <div v-if="showThree" ref="threeContainer" style="position:fixed; top:0; left:0; width:100vw; height:100vh; z-index:9999; background:black;" /> -->
	<!-- ⬇️ 使用元件版 -->
	<ThreePenguinShow />

    <div
      v-if="
        ['dashboard', 'mapview'].includes(authStore.currentPath) &&
          !authStore.isMobile &&
          !authStore.isNarrowDevice
      "
      class="app-update"
    >
      <p>下次更新：{{ formattedTimeToUpdate }}</p>
    </div>
  </div>
</template>

<style scoped lang="scss">
.app {
	&-container {
		max-width: 100vw;
		max-height: 100vh;
		max-height: calc(var(--vh) * 100);
	}

	&-content {
		width: 100vw;
		max-width: 100vw;
		height: calc(100vh - 60px);
		height: calc(var(--vh) * 100 - 60px);
		display: flex;

		&-main {
			width: 100%;
			display: flex;
			flex-direction: column;
		}
	}

	&-update {
		position: fixed;
		bottom: 0;
		right: 20px;
		color: white;
		opacity: 0.3;
		transition: opacity 0.3s;
		user-select: none;

		p {
			color: var(--color-complement-text);
		}

		&:hover {
			opacity: 1;
		}
	}
}
</style>
