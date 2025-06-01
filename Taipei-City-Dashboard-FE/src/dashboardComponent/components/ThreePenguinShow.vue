<template>
  <div ref="threeContainer" class="three-canvas" />
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, nextTick } from 'vue'
import * as THREE from 'three'
import { GLTFLoader } from 'three/examples/jsm/loaders/GLTFLoader.js'
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls.js'

const threeContainer = ref(null)
let renderer, scene, camera, model, controls, animationId
const cleanCubes = [] // 清潔方塊（綠色）
const dirtyItems = [] // 菸蒂飛行物件
let cigaretteModel = null // 預載菸蒂模型

onMounted(() => {
  nextTick(() => {
    if (!threeContainer.value) return
	// 載入水瓶模型
	const waterBottleLoader = new GLTFLoader()
	waterBottleLoader.load('/models/Water_bottle.glb', (gltf) => {
	waterBottleModel = gltf.scene
	})

    // 場景與相機
    scene = new THREE.Scene()
    scene.background = null
    camera = new THREE.PerspectiveCamera(75, 1, 0.1, 1000)
    camera.position.set(0, 2, 5)

    // 渲染器
    renderer = new THREE.WebGLRenderer({ antialias: true, alpha: true })
    renderer.setSize(300, 300)
    renderer.shadowMap.enabled = true
    renderer.shadowMap.type = THREE.PCFSoftShadowMap
    threeContainer.value.appendChild(renderer.domElement)

    // 控制器
    controls = new OrbitControls(camera, renderer.domElement)
    controls.enableDamping = true
    controls.dampingFactor = 0.05
    controls.enableZoom = false
    controls.enablePan = false

    // 燈光
    scene.add(new THREE.AmbientLight(0xffffff, 0.5))
    const dirLight = new THREE.DirectionalLight(0xffffff, 1)
    dirLight.position.set(5, 10, 7)
    dirLight.castShadow = true
    dirLight.shadow.mapSize.set(60, 60)
    dirLight.shadow.bias = -0.0005
    scene.add(dirLight)

    // 地板
    const plane = new THREE.Mesh(
      new THREE.PlaneGeometry(500, 500),
      new THREE.ShadowMaterial({ opacity: 0.3 })
    )
    plane.rotation.x = -Math.PI / 2
    plane.position.y = -1
    plane.receiveShadow = true
    scene.add(plane)

    // 載入企鵝模型
    const loader = new GLTFLoader()
    loader.load('/models/Codefest_Penguin_0531155635_texture.glb', (gltf) => {
      model = gltf.scene
      model.scale.set(0.9, 0.9, 0.9 )
      model.traverse(obj => {
        if (obj.isMesh) {
          obj.castShadow = true
          obj.receiveShadow = true
        }
      })
      scene.add(model)
      animate()
    })

    // 載入菸蒂模型
    const cigaretteLoader = new GLTFLoader()
    cigaretteLoader.load('/models/Cigarettes.glb', (gltf) => {
      cigaretteModel = gltf.scene
    })

    // 左鍵丟髒（菸蒂）、右鍵丟乾淨球
    window.addEventListener('click', throwCigarette)
    window.addEventListener('contextmenu', (e) => {
      e.preventDefault()
      throwCleanerToPenguin()
    })

    // R 鍵重設顏色
    window.addEventListener('keydown', (e) => {
      if (e.key === 'r' && model) {
        model.traverse(obj => {
          if (obj.isMesh && obj.material?.color) {
            obj.material.color.set(0xffffff)
          }
        })
      }
    })
  })
})

let waterBottleModel = null // 預載水瓶模型

function throwCleanerToPenguin() {
  if (!model || !waterBottleModel) return

  const bottle = waterBottleModel.clone()
  bottle.scale.set(1.4, 1.4, 1.4)
  bottle.position.copy(camera.position)

  // 可選：讓水瓶變亮藍
  bottle.traverse(child => {
    if (child.isMesh && child.material) {
      child.material.color.set(0x66ccff) // 淺藍水感
      child.material.emissive = new THREE.Color(0x99ccff)
      child.material.emissiveIntensity = 0.6
    }
  })

  scene.add(bottle)

  const penguinPos = new THREE.Vector3()
  model.getWorldPosition(penguinPos)
  const direction = new THREE.Vector3().subVectors(penguinPos, camera.position).normalize()

  cleanCubes.push({
    mesh: bottle,
    dir: direction,
    rotSpeed: new THREE.Vector3(Math.random(), Math.random(), Math.random())
  })
}


// 丟菸蒂（左鍵）
// 丟菸蒂（左鍵）
function throwCigarette() {
  if (!model || !cigaretteModel) return

  const cig = cigaretteModel.clone()
  cig.scale.set(1.6, 1.6, 1.6)
  cig.position.copy(camera.position)

  // 模型變白 + 發光
  cig.traverse(child => {
    if (child.isMesh && child.material) {
      // 強制換成純白材質
      child.material.color.set(0xffffff)
    //   child.material.emissive = new THREE.Color(0xff3300) // 發橘光
      child.material.emissiveIntensity = 1
    }
  })

  // 火球（火光效果）
//   const fire = new THREE.Mesh(
//     new THREE.SphereGeometry(0.1, 16, 16),
//     new THREE.MeshBasicMaterial({ color: 0xff3300 })
//   )
//   fire.position.set(0, 0, 0.2)
//   cig.add(fire)

  scene.add(cig)

  const penguinPos = new THREE.Vector3()
  model.getWorldPosition(penguinPos)
  const direction = new THREE.Vector3().subVectors(penguinPos, camera.position).normalize()

  dirtyItems.push({
    mesh: cig,
    dir: direction,
    rotSpeed: new THREE.Vector3(Math.random(), Math.random(), Math.random())
  })
}


function animate() {
  animationId = requestAnimationFrame(animate)

  if (model) model.rotation.y += 0.01
  controls.update()

  // 髒物（菸蒂）更新
  for (let i = dirtyItems.length - 1; i >= 0; i--) {
    const obj = dirtyItems[i]
    obj.mesh.position.addScaledVector(obj.dir, 0.15)
	obj.mesh.rotation.x += obj.rotSpeed.x * 0.1
	obj.mesh.rotation.y += obj.rotSpeed.y * 0.1
	obj.mesh.rotation.z += obj.rotSpeed.z * 0.1


    if (model && obj.mesh.position.distanceTo(model.position) < 1) {
      model.traverse(child => {
        if (child.isMesh && child.material?.color) {
          child.material.color.multiplyScalar(0.9) // 變髒
        }
      })
      scene.remove(obj.mesh)
      dirtyItems.splice(i, 1)
      continue
    }

    if (obj.mesh.position.distanceTo(camera.position) > 10) {
      scene.remove(obj.mesh)
      dirtyItems.splice(i, 1)
    }
  }

  // 清潔球更新
  for (let i = cleanCubes.length - 1; i >= 0; i--) {
    const obj = cleanCubes[i]
    obj.mesh.position.addScaledVector(obj.dir, 0.2)

    if (model && obj.mesh.position.distanceTo(model.position) < 1) {
      model.traverse(child => {
        if (child.isMesh && child.material?.color) {
          child.material.color.multiplyScalar(1.15) // 變乾淨
        }
      })
      scene.remove(obj.mesh)
      cleanCubes.splice(i, 1)
      continue
    }

    if (obj.mesh.position.distanceTo(camera.position) > 10) {
      scene.remove(obj.mesh)
      cleanCubes.splice(i, 1)
    }
  }

  renderer.render(scene, camera)
}

onBeforeUnmount(() => {
  cancelAnimationFrame(animationId)
  renderer?.dispose?.()
  window.removeEventListener('click', throwCigarette)
  window.removeEventListener('contextmenu', throwCleanerToPenguin)
})
</script>




<style scoped>
.three-canvas {
  position: fixed;
  bottom: 75vh;
  right: 75vw;
  width: 300px;
  height: 300px;
  z-index: 1;
  background: transparent;
  /* pointer-events: none; */
}
canvas {
  display: block;
}
</style>
