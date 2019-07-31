package mail

const (
	MSG_LINE_HEIGHT = 17
	MSG_ERR         = `%s<pre class="info">    系统发生严重错误!</pre>
               <div class="panel panel-default">
                   <div class="panel-heading-danger">堆栈信息</div>
                       <div class="panel-body">
                           <textarea class="msg-text" style="height:%spx">%s</textarea>
                       </div>
                   <div class="panel-footer"></div>
               </div>
               <div class="panel panel-default">
                     <div class="panel-heading-success">错误描述</div>
                         <div class="panel-body">
                             <textarea class="msg-text" style="height:100px">%s</textarea>
                         </div>
                     <div class="panel-footer"></div>
               </div>`

	MSG_CSS = `<style type="text/css">
              .panel{ margin-bottom: 20px; background-color: #fff; border: 1px solid transparent; border-radius: 4px; box-shadow: 0 1px 1px rgba(0,0,0,.05); margin-left: 20px; margin-right: 20px; }
              .panel-body{ padding-left: 5px; padding-right: 5px; }
              .panel-heading{ padding: 10px 15px; border-bottom: 1px solid transparent; border-top-left-radius: 3px; border-top-right-radius: 3px; color: #31708f; background-color: #d9edf7; border-color: #bce8f1; }
              .panel-heading-danger{ padding: 10px 15px; border-bottom: 1px solid transparent; border-top-left-radius: 3px; border-top-right-radius: 3px; color: #a94442; background-color: #f2dede; border-color: #ebccd1; }
              .panel-heading-success{ padding: 10px 15px; border-bottom: 1px solid transparent; border-top-left-radius: 3px; border-top-right-radius: 3px; color: #3c763d; background-color: #dff0d8; border-color: #d6e9c6; }
              .panel-heading > .dropdown .dropdown-toggle{ color: inherit; }
              .panel-title{ margin-top: 0; margin-bottom: 0; font-size: 16px; color: inherit; }
              .panel-title > a{ color: inherit; }
              .panel-footer{ padding:10px 15px; background-color: #f5f5f5; border-top: 1px solid #ddd; border-bottom-right-radius: 3px; border-bottom-left-radius: 3px; }
              .panel-default{ border-color: #bce8f1; }
              .panel-default > .panel-heading{ color: #31708f; background-color: #d9edf7; border-color: #bce8f1; }
              .panel-default > .panel-heading + .panel-collapse .panel-body{ border-top-color: #ddd; }
              .panel-default > .panel-footer + .panel-collapse .panel-body{ border-bottom-color: #ddd; }
              .msg-text{ width: 100%; border: 0px solid; }
              .info{ font-size: 14px; color: red; }</style>`
)
