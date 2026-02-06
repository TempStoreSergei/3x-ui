<script setup lang="ts">
import { onMounted } from 'vue'
import {
  SettingOutlined,
  SafetyOutlined,
  MessageOutlined,
  LinkOutlined,
  CloudOutlined,
  SaveOutlined,
  RobotOutlined,
} from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { useSettingStore } from '@/stores/setting'

const settingStore = useSettingStore()

async function handleSave() {
  const res = await settingStore.saveSettings()
  if (res?.success) {
    message.success('Settings saved successfully')
  } else {
    message.error(res?.msg || 'Failed to save settings')
  }
}

function openBotFather() {
  window.open('https://t.me/BotFather', '_blank')
}

function openTgBot() {
  const token = settingStore.allSetting?.tgBotToken
  if (!token) {
    message.warning('Please enter a bot token first')
    return
  }
  fetch(`https://api.telegram.org/bot${token}/getMe`)
    .then((r) => r.json())
    .then((data: any) => {
      if (data.ok && data.result?.username) {
        window.open(`https://t.me/${data.result.username}`, '_blank')
      } else {
        message.error('Invalid bot token')
      }
    })
    .catch(() => message.error('Could not connect to Telegram API'))
}

onMounted(() => {
  settingStore.fetchSettings()
})
</script>

<template>
  <a-spin :spinning="settingStore.loading && !settingStore.allSetting">
    <template v-if="settingStore.allSetting">
      <a-card size="small" style="margin-bottom: 16px">
        <a-button type="primary" @click="handleSave" :loading="settingStore.loading">
          <SaveOutlined /> Save Settings
        </a-button>
      </a-card>

      <a-card>
        <a-tabs>
          <!-- General -->
          <a-tab-pane key="general">
            <template #tab><SettingOutlined /> General</template>
            <a-form layout="vertical" style="max-width: 600px">
              <a-form-item label="Listen IP">
                <a-input v-model:value="settingStore.allSetting.webListen" placeholder="0.0.0.0 or empty" />
              </a-form-item>
              <a-form-item label="Port">
                <a-input-number v-model:value="settingStore.allSetting.webPort" :min="1" :max="65535" style="width: 100%" />
              </a-form-item>
              <a-form-item label="Base Path">
                <a-input v-model:value="settingStore.allSetting.webBasePath" placeholder="/" />
              </a-form-item>
              <a-form-item label="Certificate File">
                <a-input v-model:value="settingStore.allSetting.webCertFile" placeholder="/path/to/cert.pem" />
              </a-form-item>
              <a-form-item label="Key File">
                <a-input v-model:value="settingStore.allSetting.webKeyFile" placeholder="/path/to/key.pem" />
              </a-form-item>
              <a-form-item label="Session Max Age (minutes)">
                <a-input-number v-model:value="settingStore.allSetting.sessionMaxAge" :min="1" style="width: 100%" />
              </a-form-item>
              <a-form-item label="Time Location">
                <a-input v-model:value="settingStore.allSetting.timeLocation" placeholder="Asia/Tehran" />
              </a-form-item>
            </a-form>
          </a-tab-pane>

          <!-- Security -->
          <a-tab-pane key="security">
            <template #tab><SafetyOutlined /> Security</template>
            <a-form layout="vertical" style="max-width: 600px">
              <a-form-item label="Web Domain">
                <a-input v-model:value="settingStore.allSetting.webDomain" placeholder="example.com" />
              </a-form-item>
              <a-form-item label="Page Size">
                <a-input-number v-model:value="settingStore.allSetting.pageSize" :min="0" style="width: 100%" />
              </a-form-item>
              <a-form-item label="Expire Diff (days)">
                <a-input-number v-model:value="settingStore.allSetting.expireDiff" :min="0" style="width: 100%" />
              </a-form-item>
              <a-form-item label="Traffic Diff (GB)">
                <a-input-number v-model:value="settingStore.allSetting.trafficDiff" :min="0" style="width: 100%" />
              </a-form-item>
            </a-form>
          </a-tab-pane>

          <!-- Telegram -->
          <a-tab-pane key="telegram">
            <template #tab><MessageOutlined /> Telegram</template>
            <a-form layout="vertical" style="max-width: 600px">
              <a-card size="small" style="margin-bottom: 16px; background: #f0f5ff; border-radius: 8px">
                <h4 style="margin: 0 0 8px"><RobotOutlined /> Quick Connect</h4>
                <p style="font-size: 12px; color: #666">
                  1. Create a bot via @BotFather in Telegram<br />
                  2. Paste the token below<br />
                  3. Click "Open Bot" to start chatting
                </p>
                <a-space>
                  <a-button type="primary" @click="openBotFather">
                    <RobotOutlined /> Open @BotFather
                  </a-button>
                  <a-button @click="openTgBot" :disabled="!settingStore.allSetting.tgBotToken">
                    <MessageOutlined /> Open My Bot
                  </a-button>
                </a-space>
              </a-card>
              <a-form-item label="Enable Bot">
                <a-switch v-model:checked="settingStore.allSetting.tgBotEnable" />
              </a-form-item>
              <a-form-item label="Bot Token">
                <a-input v-model:value="settingStore.allSetting.tgBotToken" placeholder="123456:ABC-DEF..." />
              </a-form-item>
              <a-form-item label="Chat ID">
                <a-input v-model:value="settingStore.allSetting.tgBotChatId" placeholder="Your chat ID" />
              </a-form-item>
              <a-form-item label="Notification Time">
                <a-input v-model:value="settingStore.allSetting.tgRunTime" placeholder="@daily" />
              </a-form-item>
              <a-form-item label="Enable Backup Notifications">
                <a-switch v-model:checked="settingStore.allSetting.tgBotBackup" />
              </a-form-item>
              <a-form-item label="Enable Login Notifications">
                <a-switch v-model:checked="settingStore.allSetting.tgBotLoginNotify" />
              </a-form-item>
              <a-form-item label="CPU Alert Threshold (%)">
                <a-input-number v-model:value="settingStore.allSetting.tgCpu" :min="0" :max="100" style="width: 100%" />
              </a-form-item>
            </a-form>
          </a-tab-pane>

          <!-- Subscription -->
          <a-tab-pane key="subscription">
            <template #tab><LinkOutlined /> Subscription</template>
            <a-form layout="vertical" style="max-width: 600px">
              <a-form-item label="Enable Subscription">
                <a-switch v-model:checked="settingStore.allSetting.subEnable" />
              </a-form-item>
              <a-form-item label="Enable JSON Subscription">
                <a-switch v-model:checked="settingStore.allSetting.subJsonEnable" />
              </a-form-item>
              <a-form-item label="Listen IP">
                <a-input v-model:value="settingStore.allSetting.subListen" />
              </a-form-item>
              <a-form-item label="Port">
                <a-input-number v-model:value="settingStore.allSetting.subPort" :min="1" :max="65535" style="width: 100%" />
              </a-form-item>
              <a-form-item label="Domain">
                <a-input v-model:value="settingStore.allSetting.subDomain" placeholder="sub.example.com" />
              </a-form-item>
              <a-form-item label="Path">
                <a-input v-model:value="settingStore.allSetting.subPath" placeholder="/sub/" />
              </a-form-item>
              <a-form-item label="Subscription URI">
                <a-input v-model:value="settingStore.allSetting.subURI" placeholder="https://sub.example.com/sub/" />
              </a-form-item>
              <a-form-item label="Update Interval (hours)">
                <a-input-number v-model:value="settingStore.allSetting.subUpdates" :min="1" style="width: 100%" />
                <p style="font-size: 12px; color: #888; margin-top: 4px">
                  Recommended: 60 min for rotation setups, 1440 min (24h) for stability.
                </p>
              </a-form-item>
              <a-form-item label="Title">
                <a-input v-model:value="settingStore.allSetting.subTitle" />
              </a-form-item>
              <a-form-item label="Encrypt Subscription">
                <a-switch v-model:checked="settingStore.allSetting.subEncrypt" />
              </a-form-item>
              <a-form-item label="Show Info">
                <a-switch v-model:checked="settingStore.allSetting.subShowInfo" />
              </a-form-item>
            </a-form>
          </a-tab-pane>

          <!-- Cloudflare & DNS -->
          <a-tab-pane key="cloudflare">
            <template #tab><CloudOutlined /> Cloudflare &amp; DNS</template>
            <a-form layout="vertical" style="max-width: 600px">
              <a-card size="small" style="margin-bottom: 16px">
                <h4 style="margin: 0 0 8px">Cloudflare API</h4>
                <a-form-item label="API Token">
                  <a-input-password v-model:value="settingStore.allSetting.cloudflareAPIToken" placeholder="Enter Cloudflare API token" />
                </a-form-item>
                <a-form-item label="Default Zone ID">
                  <a-input v-model:value="settingStore.allSetting.cloudflareZoneID" placeholder="Zone ID (optional)" />
                </a-form-item>
              </a-card>
              <a-card size="small" style="margin-bottom: 16px">
                <h4 style="margin: 0 0 8px">NextDNS</h4>
                <a-form-item label="Profile ID">
                  <a-input v-model:value="settingStore.allSetting.nextDnsProfileId" placeholder="e.g. abc123" />
                  <p style="font-size: 12px; color: #888; margin-top: 4px">
                    Get your profile ID from <a href="https://my.nextdns.io" target="_blank">my.nextdns.io</a>
                  </p>
                </a-form-item>
              </a-card>
              <a-alert type="info" show-icon message="How to use" style="margin-top: 8px">
                <template #description>
                  1. Get a Cloudflare API Token (DNS Edit permission)<br />
                  2. Paste it above and save settings<br />
                  3. For NextDNS: create a profile at my.nextdns.io<br />
                  4. Run Quick Setup from Inbounds page to auto-configure
                </template>
              </a-alert>
            </a-form>
          </a-tab-pane>
        </a-tabs>
      </a-card>
    </template>
  </a-spin>
</template>
