import pexpect
import json
import os

from lib.scanner import config


# 子域收集函数
def oneforall_scan(target):
    cmd = "python3 %s --target %s --fmt json run" % (config.oneforall_exe_path,
                                                     target)
    child = pexpect.spawn(cmd, timeout=None)
    print("正在收集子域名:", target)
    # 如果存在结果文件，跳过子域名收集
    result_path="%s/%s.json"%(config.oneforall_result_path,target) 
    if not os.path.exists(result_path): 
        # 进行子域名收集
        child.expect("The subdomain result for %s:" % target)
        result_path = child.readline().decode("utf-8").strip('\x1b[0m\r\n').strip()
        child.expect(pexpect.EOF)
        print("子域名收集结束:", result_path)
    with open(result_path) as result_file:
        results = json.load(result_file)
        return results


if __name__ == "__main__":
    oneforall_scan("motorola.com.cn")