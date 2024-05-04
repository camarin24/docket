<script lang="ts" setup>
import { NH1, NSelect, NTable, NDescriptions, NDescriptionsItem, NText, NIcon, NButton, NCard, NDropdown, NTabs, NTabPane } from 'naive-ui';
import { Search, File, Calendar, Folder, List, ArrowLeft, LayoutGrid } from 'lucide-vue-next';
import { GetRecentDocuments, QueryDocuments } from "../../wailsjs/go/main/App"
import { h, onMounted, ref } from 'vue'
const options = [
    {
        label: 'Marina Bay Sands',
        key: 'marina bay sands',
    },
    {
        label: "Brown's Hotel, London",
        key: "brown's hotel, london"
    },
    {
        label: 'Atlantis Bahamas, Nassau',
        key: 'atlantis nahamas, nassau'
    },
    {
        label: 'The Beverly Hills Hotel, Los Angeles',
        key: 'the beverly hills hotel, los angeles'
    }
]


const documents = ref<any[]>([])
const queryDocuments = ref<any[]>([])
const loading = ref(false)
const selectedDocument = ref<null | any>(null)

onMounted(async () => {
    const data = await GetRecentDocuments()
    documents.value = data;
})


const onSearch = async (query: string) => {
    if (!query.length) {
        queryDocuments.value = []
        return
    }
    loading.value = true
    const result = await QueryDocuments(query)
    if (result) {
        queryDocuments.value = result.map((v) => {
            return {
                label: v.name,
                value: v
            }
        })
    }
    loading.value = false
}

const getContent = () => { 
    return selectedDocument.value.content.replaceAll("\\n", "<br/>")
}

</script>

<template>
    <div v-if="!selectedDocument">
        <div class="flex-col flex items-center justify-center pt-16 pb-10">
            <n-h1>
                <n-text type="primary" strong>
                    Bienvenido a Docket
                </n-text>
            </n-h1>
            <n-select v-model:value="selectedDocument" :loading="loading" :options="queryDocuments" remote clearable
                filterable @search="onSearch" size="large" placeholder="Buscar en Docket" class="!w-1/2">
                <template #suffix>
                    <n-text code strong>
                        ⌘ + k
                    </n-text>
                </template>
                <template #prefix>
                    <n-icon :component="Search" />
                </template>
            </n-select>
            <div class="flex space-x-5 justify-center mt-5">
                <n-dropdown trigger="hover" :options="options" :placement="'bottom-start'">
                    <n-button>
                        Tipo
                        <template #icon>
                            <n-icon :component="File" />
                        </template>
                    </n-button>
                </n-dropdown>
                <n-dropdown trigger="hover" :options="options" :placement="'bottom-start'">
                    <n-button>
                        Modificado
                        <template #icon>
                            <n-icon :component="Calendar" />
                        </template>
                    </n-button>
                </n-dropdown>
                <n-dropdown trigger="hover" :options="options" :placement="'bottom-start'">
                    <n-button>
                        Ubicación
                        <template #icon>
                            <n-icon :component="Folder" />
                        </template>
                    </n-button>
                </n-dropdown>
            </div>
        </div>
        <div class="flex justify-end">
            <n-tabs class="w-40 !p-0" tab-class="p-0" type="segment" animated size="small" :bar-width="50">
                <n-tab-pane name="list" class="!p-0" :tab="h(List, { width: 20, height: 20 })">
                </n-tab-pane>
                <n-tab-pane name="grid" :tab="h(LayoutGrid, { width: 20, height: 20 })">
                </n-tab-pane>
            </n-tabs>
        </div>
        <n-table :bordered="false" :single-line="true">
            <thead>
                <tr>
                    <th>Nombre</th>
                    <th>Añadido</th>
                    <th>Tamaño</th>
                    <th></th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="d in documents" :key="d.id" @click="selectedDocument = d">
                    <td>{{ d.name }}</td>
                    <td>{{ d.CreatedAt }}</td>
                    <td>{{ d.size }} bytes</td>
                    <td>...</td>
                </tr>
            </tbody>
        </n-table>
    </div>
    <n-card v-if="selectedDocument">
        <div class="flex space-x-4 items-center">
            <n-button @click="selectedDocument = null">
                <n-icon :component="ArrowLeft" />
            </n-button>
            <div class="w-full">
                <n-descriptions label-placement="left" size="small" :column="4" :title="selectedDocument.name">
                <n-descriptions-item label="Fecha de escaneo">
                    {{ selectedDocument.CreatedAt }}
                </n-descriptions-item>
                <n-descriptions-item label="Modificado">
                    {{ selectedDocument.UpdatedAt }}
                </n-descriptions-item>
                <n-descriptions-item label="Ubicación">
                    {{ selectedDocument.originalPath }}
                </n-descriptions-item>
                <n-descriptions-item label="Tamaño">
                    {{ selectedDocument.size }} bytes
                </n-descriptions-item>
            </n-descriptions>
            </div>
        </div>
    </n-card>
    <div v-if="selectedDocument" class="flex">
        <n-card>
        <iframe src="../assets/file.pdf" frameborder="0"></iframe>
            <p v-html="getContent()"></p>
        </n-card>
        <n-card class="w-1/3">
            <n-descriptions label-placement="left" title="Metadata" :column="1">
                <n-descriptions-item v-for="(v, k) in JSON.parse(selectedDocument.metaData)" :key="k" :label="k as any">
                    {{ v }}
                </n-descriptions-item>
            </n-descriptions>
        </n-card>
    </div>
</template>