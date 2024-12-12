package utils

func CheckPassword(pwd string) (sucess bool, pwdTotalCount int64) {

	r, err := Scalar("select count(1) from Passwords")
	if err != nil && r == nil {
		return false, 0
	}
	val1, ok := r.(int64)
	if ok {
		if val1 == 0 {
			return true, val1
		} else {

			r, err = Scalar("select count(1) from Passwords where pwd=?", pwd)
			if err != nil && r == nil {
				return false, val1
			}
			val2, ok := r.(int64)
			if ok {
				return val2 > 0, val1
			}
			return false, val1
		}
	} else {
		return false, 0
	}
}

func EditePassword(oldPwd string, newPwd string) (success bool, msg string) {
	sucess, count := CheckPassword(oldPwd)
	if sucess {
		r, err := ExecNonQuery("insert into Passwords (pwd)values(?)", newPwd)
		if r < 1 || err != nil {
			return false, "修改失败1" + err.Error()
		}
		if count > 0 {

			r, err = ExecNonQuery("delete from Passwords where pwd=?", oldPwd)
			if r < 1 || err != nil {
				return false, "修改失败2"
			} else {
				return true, "修改成功"
			}
		} else {
			return true, "修改成功"

		}
	} else {
		return false, "旧密码错误"
	}
}
