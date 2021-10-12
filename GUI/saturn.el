;;; saturn.el --- Emacs interface for saturn, the animesaturn.it scraper  -*- lexical-binding: t; -*-

;; Copyright (C) 2021  Omar Polo

;; Author: Omar Polo <op@omarpolo.com>
;; Keywords: multimedia

;; This program is free software; you can redistribute it and/or modify
;; it under the terms of the GNU General Public License as published by
;; the Free Software Foundation, either version 3 of the License, or
;; (at your option) any later version.

;; This program is distributed in the hope that it will be useful,
;; but WITHOUT ANY WARRANTY; without even the implied warranty of
;; MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
;; GNU General Public License for more details.

;; You should have received a copy of the GNU General Public License
;; along with this program.  If not, see <https://www.gnu.org/licenses/>.

;;; Commentary:

;; saturn.el is an Emacs interface for saturn(1), the animesaturn.it
;; scraper.  It allows to search animes and download or stream
;; multiple episodes using a dired-like interface.

;;; Code:
(require 'cl-lib)
(require 'ewoc)

(defgroup saturn nil
  "Saturn."
  :group 'saturn)

(defcustom saturn-cmd "saturn"
  "Path to the saturn executable."
  :type 'string)

(defun saturn--completing-read-anime ()
  "Search and select an anime."
  (with-temp-buffer
    (let* ((anime (read-from-minibuffer "Which anime: ")))
      (message "searching for \"%s\"..." anime)
      (unless (zerop (call-process saturn-cmd nil t nil "-search" anime))
        (error "Saturn failed: %s" (buffer-string)))
      (completing-read "Select the anime: "
                       (split-string (buffer-string) "\n" t " \t")))))

(defvar saturn--current-anime nil
  "Anime associated with the *saturn* buffer.")

(defvar saturn--ewoc nil
  "Ewoc instance.")

(defun saturn--pp (data)
  "Pretty print DATA for ewoc."
  (cl-destructuring-bind (marked-p id ep url) data
    (insert
     (if marked-p "*" " ") " "
     id "\t"
     ep "\t"
     url)))

(defun saturn--parse-episode (line)
  "Parse LINE into a list (id ep-name url)."
  (cl-destructuring-bind (id url) (split-string line)
    (list (substring id 3)
          (file-name-nondirectory url)
          url)))

(defun saturn--search-episodes (anime)
  "Return a list of episodes for the given ANIME."
  (with-temp-buffer
    (unless (zerop (call-process saturn-cmd nil t nil
                                 "-fetch" anime))
      (error "Saturn failed: %s" (buffer-string)))
    (mapcar #'saturn--parse-episode
            (split-string (buffer-string) "\n" t " \t"))))

(defun saturn-plot (anime)
  "Return the plot for the given ANIME."
  (with-temp-buffer
    (when (zerop (call-process saturn-cmd nil t nil "-plot" anime))
      (let ((plot (string-trim (buffer-string))))
        (unless (string= plot "")
          plot)))))

;;;###autoload
(defun saturn (anime)
  "List episodes of the given ANIME."
  (interactive (list (saturn--completing-read-anime)))
  (if-let (buf (get-buffer "*saturn*"))
      (kill-buffer buf))
  (switch-to-buffer (get-buffer-create "*saturn*"))
  (saturn-mode)
  (let ((plot (saturn-plot anime)))
    (setq-local saturn--current-anime anime
                saturn--ewoc (ewoc-create #'saturn--pp
                                          (or plot
                                              (concat "Results for " anime ":\n")))))
  (dolist (ep (saturn--search-episodes anime))
    (when ep
      (ewoc-enter-last saturn--ewoc (cons nil ep)))))

(defvar saturn-mode-map
  (let ((m (make-sparse-keymap)))
    (define-key m (kbd "n") #'next-line)
    (define-key m (kbd "p") #'previous-line)
    (define-key m (kbd "m") #'saturn-mark)
    (define-key m (kbd "u") #'saturn-unmark)
    (define-key m (kbd "d") #'saturn-download)
    (define-key m (kbd "s") #'saturn-stream)
    m))

(defun saturn--revert (&rest _)
  "Revert the *saturn* buffer."
  (call-interactively #'saturn))

(defvar saturn-mode-hook (list #'visual-line-mode)
  "Hook for `saturn-mode'.")

(define-derived-mode saturn-mode special-mode "Saturn"
  "Major mode for Saturn episode list."
  (setq-local revert-buffer-function #'saturn--revert))

(defun saturn--set-at-point (val)
  "Mark/unmork episode at point, depending on VAL."
  (if-let (e (ewoc-locate saturn--ewoc))
    (let* ((data (ewoc-data e)))
      (setf (car data) val)
      (ewoc-set-data e data)
      (ewoc-invalidate saturn--ewoc e)
      (forward-line 1))
    (error "No episode at point")))

(defun saturn-mark ()
  "Mark the episode at point."
  (interactive "" saturn-mode)
  (saturn--set-at-point t))

(defun saturn-unmark ()
  "Unmark the episode at point."
  (interactive "" saturn-mode)
  (saturn--set-at-point nil))

(defun saturn--marked ()
  "Return a list of marked episodes ids or the one at point."
  (or (mapcar #'cadr (ewoc-collect saturn--ewoc (lambda (data) (car data))))
      (list (cadr (ewoc-data (ewoc-locate saturn--ewoc))))))



(defvar saturn--stream-proc nil
  "The streaming process.")

(defun saturn-stream (anime episodes)
  "Stream EPISODES of the given ANIME."
  (interactive (list saturn--current-anime
                     (saturn--marked))
               saturn-mode)
  (when (process-live-p saturn--stream-proc)
    (error "[saturn] the streaming is already running"))
  (let* ((process-connection-type t)
         (proc (start-process "saturn stream"
                              "*saturn-stream*"
                              saturn-cmd
                              "-fetch" anime
                              "-stream" (mapconcat #'identity
                                                   episodes ","))))
    (message "Streaming episode%s %s"
             (if (cdr episodes) "s" "") ; poor man' pluralizing
             (mapconcat #'identity episodes ", "))
    (setq saturn--stream-proc proc)))



(defvar saturn-download-mode-map
  (let ((m (make-sparse-keymap)))
    (define-key m (kbd "q") #'saturn-kill-download)
    m))

(define-derived-mode saturn-download-mode special-mode "saturn:download"
  "Major mode used for downloading episodes in `saturn'.")

(defvar saturn--downlad-proc nil
  "The download process.")

(defun saturn-kill-download ()
  "Kill the download process."
  (interactive)
  (when (process-live-p saturn--downlad-proc)
    (when (y-or-n-p "Really kill the saturn process? ")
      (kill-process saturn--downlad-proc))))

(defun saturn--download-window (buf height)
  "Prepare a window for BUF with the given HEIGHT."
  (let* ((root (frame-root-window))
         (win (split-window root (- (window-height root)
                                    height))))
    (set-window-buffer win buf)
    (set-window-dedicated-p win t)
    win))

(defun saturn--download-proc-sentinel (proc _msg)
  "Sentinel function for the download PROC."
  (if (and (eq (process-status proc) 'exit)
           (zerop (process-exit-status proc)))
      (let ((buf (process-buffer saturn--downlad-proc)))
        (mapc #'delete-window (get-buffer-window-list buf))
        (kill-buffer buf))
    (message "[saturn] Fetch failed")))

(defun saturn--download-proc-filter (proc msg)
  "Filter function.
PROC is the process and MSG a chunk of its output."
  (when-let (buf (process-buffer proc))
    (with-current-buffer buf
      (let ((inhibit-read-only t))
        (goto-char (point-max))
        (if (not (string-match ".*\r\\(.*\\)" msg))
            (insert msg)
          ;; handle \r
          (end-of-line)
          (let ((end (point)))
            (beginning-of-line)
            (delete-region (point) end))
          (insert (match-string 1 msg)))))))

(defun saturn-download (anime episodes)
  "Download EPISODES of the given ANIME."
  (interactive (list saturn--current-anime
                     (saturn--marked))
               saturn-mode)
  (when (process-live-p saturn--downlad-proc)
    (error "[saturn] the download is already running"))
  (saturn--do-download anime episodes))

(defun saturn--do-download (anime episodes)
  "Fetch EPISODES of the given ANIME."
  (let* ((process-connection-type t)
         (proc (start-process "saturn download"
                              "*saturn-download*"
                              saturn-cmd
                              "-fetch" anime
                              "-down"  (mapconcat #'identity episodes ",")))
         (buf (process-buffer proc))
         (win (saturn--download-window buf 5)))
    (when (window-live-p win)
      (with-selected-window win
        (erase-buffer)
        ;; mu4e-utils.el:mu4e~update-mail-and-index-real has a FIXME
        ;; saying that without this the process doesn't start.
        (insert "\n")
        (saturn-download-mode)))
    (process-put proc 'x-interactive t)
    (set-process-sentinel proc #'saturn--download-proc-sentinel)
    (set-process-filter proc #'saturn--download-proc-filter)
    (setq saturn--downlad-proc proc)))

(provide 'saturn)
;;; saturn.el ends here
