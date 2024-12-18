package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"time"
)

const (
	host     = "smtp.mail.ru"
	port     = "587"
	from     = "m.sayenko_site@mail.ru"
	username = "m.sayenko_site@mail.ru"
	password = "qQrxNpfikdUcxx2QzxkM"
)

type Feedback struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/feedback", feedbackHandler)
	http.HandleFunc("/feedback/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello test feedback"))
	})

	fmt.Println("Сервер запущен на порту 8000...")

	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func sendMail(to []string, subject, body string) error {
	start := time.Now()
	msg := []byte("From: " + from + "\r\n" + // Указываем отправителя
		"To: " + to[0] + "\r\n" + // Указываем получателя
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n" +
		body + "\r\n")
	auth := smtp.PlainAuth("", username, password, host)

	for i := 0; i < 3; i++ {
		err := smtp.SendMail(host+":"+port, auth, from, to, msg)
		if err == nil {
			log.Printf("Письмо успешно отправлено, этап занял: %s", time.Since(start))
			return nil
		}
		log.Printf("Ошибка отправки письма, попытка %d: %v", i+1, err)
		time.Sleep(time.Second * 2)
	}

	log.Printf("Не удалось отправить письмо за 3 попытки, этап занял: %s", time.Since(start))
	return fmt.Errorf("не удалось отправить письмо")
}

func feedbackHandler(w http.ResponseWriter, r *http.Request) {
	log.Fatal("feedbackHandler")
	w.Header().Set("Access-Control-Allow-Origin", "https://lawyer-pi.vercel.app")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3001")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "https://lawyer-pi.vercel.app")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3001")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusNoContent) // 204 No Content
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка при чтении тела запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var feedback Feedback
	if err := json.Unmarshal(body, &feedback); err != nil {
		http.Error(w, "Ошибка при декодировании JSON", http.StatusBadRequest)
		return
	}

	htmlBody := generateEmailBody(feedback)

	to := []string{"Saenko-kursk@mail.ru"}
	subject := "Новая запись на консультацию"

	if err := sendMail(to, subject, htmlBody); err != nil {
		http.Error(w, "Ошибка при отправке письма: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]string{"status": "success", "message": "Письмо отправлено"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func generateEmailBody(data Feedback) string {
	return fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<body style="font-family: Arial, sans-serif; background-color: #f4f4f4; margin: 0; padding: 0;">
		<div style="background-color: #ffffff; margin: 20px auto; padding: 20px; border-radius: 8px; box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1); max-width: 700px;">
			<div style="font-size: 24px; font-weight: bold; color: #333333; margin-bottom: 20px; text-align: center;">
				Новая запись на консультацию
			</div>
			<div style="font-size: 18px; color: #555555; line-height: 1.6;">
				<p>Здравствуйте!</p>
				<p>Пользователь оставил заявку на консультацию. Вот его данные:</p>
				<ul style="list-style: none; padding: 0;">
					<li><strong>ФИО:</strong> %s</li>
					<li><strong>Телефон:</strong> %s</li>
					<li><strong>Email:</strong> %s</li>
					<li><strong>Тип консультации:</strong> %s</li>
				</ul>
				
				<p style="text-align: center; font-size: 18px; color: #333; font-weight: bold; margin-bottom: 10px;">Суть обращения:</p>
				<p style="text-align: center; font-size: 18px; color: #333; margin-bottom: 10px;">%s</p>
				
				<div style="background-color: #e6f7ff; border-left: 4px solid #007acc; padding: 20px; margin-top: 30px; font-weight: bold;">
					<p>Пожалуйста, свяжитесь с ним для уточнения деталей.</p>
				</div>
			</div>
			<div style="font-size: 14px; color: #888888; text-align: center; margin-top: 20px;">
				Это письмо создано автоматически. Пожалуйста, не отвечайте на него.
			</div>
		</div>
	</body>
	</html>`, data.Name, data.Phone, data.Email, data.Type, data.Message)
}

/*
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"time"
)

const (
	host     = "smtp.gmail.com"
	port     = "587"
	from     = "sergeevnicolas20@gmail.com"
	username = "sergeevnicolas20@gmail.com"
	password = "rqrx wvjq nanc eolr"
)

type Feedback struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/feedback", feedbackHandler)

	fmt.Println("Сервер запущен на порту 80...")

	log.Fatal(http.ListenAndServe(":80", nil))
}

func sendMail(to []string, subject, body string) error {
	start := time.Now()
	msg := []byte("Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n" +
		body + "\r\n")
	auth := smtp.PlainAuth("", username, password, host)

	for i := 0; i < 3; i++ {
		err := smtp.SendMail(host+":"+port, auth, from, to, msg)
		if err == nil {
			log.Printf("Письмо успешно отправлено, этап занял: %s", time.Since(start))
			return nil
		}
		log.Printf("Ошибка отправки письма, попытка %d: %v", i+1, err)
		time.Sleep(time.Second * 2)
	}

	log.Printf("Не удалось отправить письмо за 3 попытки, этап занял: %s", time.Since(start))
	return fmt.Errorf("не удалось отправить письмо")
}

func feedbackHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "https://lawyer-pi.vercel.app")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3001")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "https://lawyer-pi.vercel.app")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusNoContent) // 204 No Content
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка при чтении тела запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var feedback Feedback
	if err := json.Unmarshal(body, &feedback); err != nil {
		http.Error(w, "Ошибка при декодировании JSON", http.StatusBadRequest)
		return
	}

	htmlBody := generateEmailBody(feedback)

	to := []string{"sergeevnicolas20@gmail.com"}
	subject := "Новая запись на консультацию"

	if err := sendMail(to, subject, htmlBody); err != nil {
		http.Error(w, "Ошибка при отправке письма: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]string{"status": "success", "message": "Письмо отправлено"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func generateEmailBody(data Feedback) string {
	return fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<style>
			body { font-family: Arial, sans-serif; background-color: #f4f4f4; margin: 0; padding: 0; }
			.email-container { background-color: #ffffff; margin: 20px auto; padding: 20px; border-radius: 8px; box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1); max-width: 700px; }
			.header { font-size: 24px; font-weight: bold; color: #333333; margin-bottom: 20px; text-align: center; }
			.content { font-size: 18px; color: #555555; line-height: 1.6; }
			.footer { font-size: 14px; color: #888888; text-align: center; margin-top: 20px; }
			.message-block {
				background-color: #e6f7ff;
				border-left: 4px solid #007acc;
				padding: 20px;
				margin-top: 30px;
				font-weight: bold;
				max-width: 100;
			}
			.message-title {
				text-align: center;
				font-size: 18px;
				color: #333;
				font-weight: bold;
				margin-bottom: 10px;
			}
			.message-text {
				text-align: center;
				font-size: 18px;
				color: #333;
				margin-bottom: 10px;
			}
		</style>
	</head>
	<body>
		<div class="email-container">
			<div class="header">
				Новая запись на консультацию
			</div>
			<div class="content">
				<p>Здравствуйте!</p>
				<p>Пользователь оставил заявку на консультацию. Вот его данные:</p>
				<ul>
					<li><strong>ФИО:</strong> %s</li>
					<li><strong>Телефон:</strong> %s</li>
					<li><strong>Email:</strong> %s</li>
					<li><strong>Тип консультации:</strong> %s</li>
				</ul>

				<p class="message-title">Суть обращения:</p>
				<p class="message-text">%s</p>

				<div class="message-block">
					<p>Пожалуйста, свяжитесь с ним для уточнения деталей.</p>
				</div>
			</div>
			<div class="footer">
				Это письмо создано автоматически. Пожалуйста, не отвечайте на него.
			</div>
		</div>
	</body>
	</html>`, data.Name, data.Phone, data.Email, data.Type, data.Message)
}
*/
