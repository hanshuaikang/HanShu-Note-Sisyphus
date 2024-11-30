# -*- coding: UTF-8 -*-
import time

from PIL import Image
from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.common.by import By
from selenium.webdriver.support.wait import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC


class GrafanaImageExporter:

    def __init__(self,
                 grafana_host,
                 username,
                 password,
                 dashboard_url=None,
                 height=2300,
                 output_file='dashboard.png',
                 collapse_sidebar_button_id="dock-menu-button",
                 login_button_css_name=".css-1b7vft8-button"):
        self.grafana_host = grafana_host
        self.username = username
        self.password = password
        self.dashboard_url = dashboard_url
        self.output_file = output_file
        self.collapse_sidebar_button_id = collapse_sidebar_button_id
        self.login_button_css_name = login_button_css_name
        self.height = height

        chrome_options = Options()
        chrome_options.add_argument("--headless")  # 启用无头模式
        chrome_options.add_argument("--no-sandbox")  # 解决一些权限问题
        chrome_options.add_argument("--disable-dev-shm-usage")  # 解决共享内存问题

        # Setup Chrome driver
        driver = webdriver.Chrome(options=chrome_options)
        self.driver = driver

    def export(self):
        self._do_login()
        print("now start export....")
        self._export_dashboard_image()

    def _do_login(self):
        self.driver.get(f'{self.grafana_host}/login')
        self.driver.maximize_window()
        # 等待页面加载完毕
        WebDriverWait(self.driver, 30).until(
            EC.presence_of_element_located((By.NAME, "user"))
        )
        # 找到并输入用户名
        username_input = self.driver.find_element(By.NAME, "user")
        username_input.send_keys(self.username)

        # 找到并输入密码
        password_input = self.driver.find_element(By.NAME, "password")
        password_input.send_keys(self.password)
        # 点击登录按钮
        login_button = self.driver.find_element(By.CSS_SELECTOR, self.login_button_css_name)
        login_button.click()
        # 等待三秒
        time.sleep(10)
        print("login success")

    def _collapse_sidebar(self):
        try:
            # 等待页面加载完毕
            print("now start collapse_sidebar....")
            WebDriverWait(self.driver, 30).until(
                EC.presence_of_element_located((By.ID, self.collapse_sidebar_button_id))
            )
            # 使用按钮的 ID 来找到并点击收起侧滑栏的按钮
            toggle_button = self.driver.find_element(By.ID, self.collapse_sidebar_button_id)
            toggle_button.click()
            print("now start collapse_sidebar success....")
        except Exception as e:
            print("now start collapse_sidebar failed....")
            raise e

    def _wait_for_panel_render(self):
        print("now start wait_for_panel_render....")
        # 找到所有具有 data-panel-id 的元素, 计算 panel 的数量
        WebDriverWait(self.driver, 30).until(
            EC.presence_of_element_located((By.CSS_SELECTOR, "[data-viz-panel-key]"))
        )

        elements = self.driver.find_elements(By.CSS_SELECTOR, '[data-viz-panel-key]')
        panel_count = len(elements)

        # # 等待所有的 panel 都加载完成
        wait_panel_count = 0
        while True:
            wait_panel_count += 1
            # 每个仪表盘等待三秒的时间加载
            time.sleep(3)
            print(f'Waiting for panel render {wait_panel_count}/{panel_count}')
            if wait_panel_count > panel_count:
                break

    def _export_dashboard_image(self):
        print("now start export_dashboard_image....")
        self.driver.set_window_size(1920, self.height)
        self.driver.get(self.dashboard_url)
        # 等待十秒等页面完全跳转到仪表盘界面
        time.sleep(10)
        # 收起侧滑栏
        self._collapse_sidebar()
        # 等待 30 秒等页面加载完成
        self._wait_for_panel_render()
        self.driver.save_screenshot(self.output_file)
        self._remove_bottom_blank()
        print("now start export_dashboard_image success....")

    def _remove_bottom_blank(self, blank_color=(244, 245, 245)):
        # 打开图像
        img = Image.open(self.output_file)
        pixels = img.load()

        # 找到最后一行非空白的像素
        bottom = img.height
        start_scan_line = max(0, img.height - 20)
        for y in range(start_scan_line - 1, -1, -1):
            for x in range(img.width):
                if pixels[x, y] != blank_color:
                    bottom = y + 1
                    break
            if bottom != img.height:
                break

        # 裁剪图像
        cropped_img = img.crop((0, 0, img.width, bottom))
        cropped_img.save(self.output_file)


if __name__ == "__main__":
    GrafanaImageExporter(
        grafana_host="xxxx",
        username="xxxx",
        password="xxx",
        height=3000,
        output_file="dashboard.png",
        dashboard_url="xxxx&theme=light",
    ).export()
